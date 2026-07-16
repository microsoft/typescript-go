package compiler_test

import (
	"context"
	"slices"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/contentmapper"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/spanmap"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

type fakeContentMapperHost struct {
	transform func(fileName string, content string) (contentmapper.Result, error)
}

func (r fakeContentMapperHost) Transform(mapper *contentmapper.Mapper, request contentmapper.Request) (contentmapper.Result, error) {
	return r.transform(request.FileName, request.Content)
}

func (fakeContentMapperHost) Close() error { return nil }

func newContentMapperProgram(t *testing.T, contentMapperHost contentmapper.Host, files map[string]string, rootFiles []string) *compiler.Program {
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
		ParsedConfig: &tsoptions.ParsedOptions{
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
		Config: config,
		Host:   compiler.NewCompilerHost("/src", fs, bundled.LibPath(), nil, nil, contentMapperHost),
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

	const transformed = "export const x = 1;\n"
	const original = "<template>x</template>\n"

	atomAll := func(origEnd int) *spanmap.SpanMap {
		return spanmap.New([]spanmap.Segment{{
			GenStart: 0, GenEnd: core.TextPos(len(transformed)),
			OrigStart: 0, OrigEnd: core.TextPos(origEnd), Kind: spanmap.KindAtom,
		}})
	}

	testCases := []struct {
		name     string
		mappings *spanmap.SpanMap
		wantCode int32
	}{
		{
			"overlap",
			spanmap.New([]spanmap.Segment{
				{GenStart: 0, GenEnd: 10, OrigStart: 0, OrigEnd: 0, Kind: spanmap.KindAtom},
				{GenStart: 5, GenEnd: core.TextPos(len(transformed)), OrigStart: 0, OrigEnd: 0, Kind: spanmap.KindAtom},
			}),
			100038,
		},
		{
			"outOfBounds",
			atomAll(len(original) + 50),
			100029,
		},
		{
			// A verbatim segment whose original text differs from the transformed text.
			"verbatimMismatch",
			spanmap.New([]spanmap.Segment{{
				GenStart: 0, GenEnd: core.TextPos(len(transformed)),
				OrigStart: 0, OrigEnd: core.TextPos(len(transformed)), Kind: spanmap.KindVerbatim,
			}}),
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
				transform: func(fileName string, content string) (contentmapper.Result, error) {
					return contentmapper.Result{
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
