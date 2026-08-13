package compiler_test

import (
	"context"
	"errors"
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

func (r fakeContentMapperHost) Refresh() error                                 { return nil }
func (r fakeContentMapperHost) Identities() ([]string, error)                  { return nil, nil }
func (r fakeContentMapperHost) Identity(*contentmapper.Mapper) (string, error) { return "test", nil }
func (r fakeContentMapperHost) WatchedFiles() ([]string, error)                { return nil, nil }
func (r fakeContentMapperHost) Close() error                                   { return nil }

func (r fakeContentMapperHost) Transform(mapper *contentmapper.Mapper, request contentmapper.Request) (contentmapper.Result, error) {
	return r.transform(request.FileName, request.Content)
}

func newContentMapperProgram(t *testing.T, contentMapperProject contentmapper.Project, files map[string]string, rootFiles []string) *compiler.Program {
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
		Host:   compiler.NewCompilerHost("/src", fs, bundled.LibPath(), nil, nil, contentMapperProject),
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
	mappings := spanmap.New([]spanmap.Segment{
		{VirtualStart: 0, VirtualEnd: 10, OriginalStart: 0, OriginalEnd: 0, Kind: spanmap.KindAtom},
		{VirtualStart: 5, VirtualEnd: core.TextPos(len(transformed)), OriginalStart: 0, OriginalEnd: 0, Kind: spanmap.KindAtom},
	})
	files := map[string]string{
		"/src/app.ts":        `import "./Component.vue";`,
		"/src/Component.vue": original,
	}
	contentMapperHost := fakeContentMapperHost{
		transform: func(fileName string, content string) (contentmapper.Result, error) {
			return contentmapper.Result{Text: transformed, VirtualExtension: ".ts", Mappings: mappings}, nil
		},
	}
	program := newContentMapperProgram(t, contentMapperHost, files, []string{"/src/app.ts"})
	diagnostics := collectContentMapperDiagnostics(program)
	found := slices.ContainsFunc(diagnostics, func(diagnostic *ast.Diagnostic) bool { return diagnostic.Code() == 100038 })
	assert.Assert(t, found, "expected an invalid mapping diagnostic, got: %v", diagnostics)
}

func TestContentMapperSourceFileState(t *testing.T) {
	t.Parallel()

	t.Run("successful synthesized empty file", func(t *testing.T) {
		t.Parallel()
		program := newContentMapperProgram(t, fakeContentMapperHost{
			transform: func(fileName string, content string) (contentmapper.Result, error) {
				return contentmapper.Result{Text: "export {};", VirtualExtension: ".ts", Mappings: spanmap.New(nil)}, nil
			},
		}, map[string]string{"/src/empty.vue": ""}, []string{"/src/empty.vue"})
		file := program.GetSourceFile("/src/empty.vue")
		assert.Assert(t, file != nil)
		assert.Equal(t, file.OriginalText(), "")
		assert.Equal(t, file.ContentMapper(), "vue-mapper@1.0.0")
		assert.Assert(t, !file.IsContentMapperFailureStub())
	})

	t.Run("failed transform", func(t *testing.T) {
		t.Parallel()
		program := newContentMapperProgram(t, fakeContentMapperHost{
			transform: func(fileName string, content string) (contentmapper.Result, error) {
				return contentmapper.Result{}, errors.New("failed")
			},
		}, map[string]string{"/src/fail.vue": "original"}, []string{"/src/fail.vue"})
		file := program.GetSourceFile("/src/fail.vue")
		assert.Assert(t, file != nil)
		assert.Equal(t, file.OriginalText(), "original")
		assert.Equal(t, file.ContentMapper(), "vue-mapper@1.0.0")
		assert.Assert(t, file.IsContentMapperFailureStub())
	})

	t.Run("project error is localized", func(t *testing.T) {
		t.Parallel()
		program := newContentMapperProgram(t, fakeContentMapperHost{
			transform: func(fileName string, content string) (contentmapper.Result, error) {
				return contentmapper.Result{}, contentmapper.NewTransformError(
					contentmapper.TransformErrorKindProject,
					&contentmapper.ProjectError{Kind: contentmapper.ProjectErrorKindMalformedResponse},
				)
			},
		}, map[string]string{"/src/fail.vue": "original"}, []string{"/src/fail.vue"})
		diagnostics := collectContentMapperDiagnostics(program)
		found := slices.ContainsFunc(diagnostics, func(diagnostic *ast.Diagnostic) bool {
			return slices.ContainsFunc(diagnostic.MessageChain(), func(message *ast.Diagnostic) bool {
				return message.Code() == 100051
			})
		})
		assert.Assert(t, found, "expected a localized project response diagnostic, got: %v", diagnostics)
	})
}
