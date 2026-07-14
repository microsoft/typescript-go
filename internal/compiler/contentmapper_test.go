package compiler_test

import (
	"context"
	"slices"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/contentmapper"
	"github.com/microsoft/typescript-go/internal/contentmapperhost"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/spanmap"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

type fakeContentMapperHost struct {
	transform func(fileName string, content string) (contentmapperhost.Result, error)
}

func (r fakeContentMapperHost) Transform(mapper *contentmapper.Mapper, request contentmapperhost.Request) (contentmapperhost.Result, error) {
	return r.transform(request.FileName, request.Content)
}

func (fakeContentMapperHost) Close() error { return nil }

func newContentMapperProgram(t *testing.T, contentMapperHost contentmapperhost.Host, files map[string]string, rootFiles []string) *compiler.Program {
	t.Helper()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}
	fs := vfstest.FromMap[any](nil, false /*useCaseSensitiveFileNames*/)
	fs = bundled.WrapFS(fs)
	for name, content := range files {
		_ = fs.WriteFile(name, content)
	}

	config := &tsoptions.ParsedCommandLine{
		ParsedConfig: &core.ParsedOptions{
			FileNames: rootFiles,
			CompilerOptions: &core.CompilerOptions{
				SkipLibCheck:     core.TSTrue,
				Module:           core.ModuleKindESNext,
				ModuleResolution: core.ModuleResolutionKindBundler,
			},
			ContentMappers: []*contentmapper.Mapper{{Definition: contentmapper.Definition{Package: "vue", Extensions: []string{".vue"}}, Manifest: contentmapper.Manifest{Name: "vue-mapper", Version: "1.0.0"}}},
		},
	}
	return compiler.NewProgram(compiler.ProgramOptions{
		Config:            config,
		Host:              compiler.NewCompilerHost("/src", fs, bundled.LibPath(), nil, nil),
		ContentMapperHost: contentMapperHost,
		// Load files on the calling goroutine for deterministic diagnostics ordering.
		SingleThreaded: core.TSTrue,
	})
}

func collectContentMapperDiagnostics(program *compiler.Program) []*ast.Diagnostic {
	ctx := context.Background()
	return slices.Concat(
		program.GetSyntacticDiagnostics(ctx, nil),
		program.GetSemanticDiagnostics(ctx, nil),
		program.GetProgramDiagnostics(),
	)
}

func TestContentMapperInvalidMappings(t *testing.T) {
	t.Parallel()

	// Mappings are a required, enforced part of the content mapper contract. Each malformed map is
	// attributed to the mapper with a specific diagnostic instead of surfacing untrustworthy positions.
	const transformed = "export const x = 1;\n"
	const original = "<template>x</template>\n"

	verbatimAll := func(kind spanmap.Kind, origEnd int) *spanmap.SpanMap {
		return spanmap.New([]spanmap.Segment{{
			GenStart: 0, GenEnd: core.TextPos(len(transformed)),
			OrigStart: 0, OrigEnd: core.TextPos(origEnd), Kind: kind,
		}})
	}

	testCases := []struct {
		name     string
		mappings *spanmap.SpanMap
		wantCode int32
	}{
		{"missing", nil, 100027},
		{
			"coverage",
			spanmap.New([]spanmap.Segment{{GenStart: 0, GenEnd: 3, OrigStart: 0, OrigEnd: 0, Kind: spanmap.KindSynthesized}}),
			100028,
		},
		{
			"outOfBounds",
			verbatimAll(spanmap.KindSynthesized, len(original)+50),
			100029,
		},
		{
			// A verbatim segment whose original text differs from the transformed text.
			"verbatimMismatch",
			verbatimAll(spanmap.KindVerbatim, len(transformed)),
			100030,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			files := map[string]string{
				"/src/app.ts":        `import "./Component.vue";`,
				"/src/Component.vue": original,
			}
			contentMapperHost := fakeContentMapperHost{
				transform: func(fileName string, content string) (contentmapperhost.Result, error) {
					return contentmapperhost.Result{
						Text:       transformed,
						ScriptKind: core.ScriptKindTS,
						Mappings:   tc.mappings,
					}, nil
				},
			}
			program := newContentMapperProgram(t, contentMapperHost, files, []string{"/src/app.ts"})
			diags := collectContentMapperDiagnostics(program)
			found := slices.ContainsFunc(diags, func(d *ast.Diagnostic) bool { return d.Code() == tc.wantCode })
			assert.Assert(t, found, "expected diagnostic TS%d attributing the invalid mapping, got: %v", tc.wantCode, diags)
		})
	}
}
