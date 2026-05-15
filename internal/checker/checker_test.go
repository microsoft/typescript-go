package checker_test

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func TestGetSymbolAtLocation(t *testing.T) {
	t.Parallel()

	content := `interface Foo {
  bar: string;
}
declare const foo: Foo;
foo.bar;`
	fs := vfstest.FromMap(map[string]string{
		"/foo.ts": content,
		"/tsconfig.json": `
				{
					"compilerOptions": {},
					"files": ["foo.ts"]
				}
			`,
	}, false /*useCaseSensitiveFileNames*/)
	fs = bundled.WrapFS(fs)

	cd := "/"
	host := compiler.NewCompilerHost(cd, fs, bundled.LibPath(), nil, nil)

	parsed, errors := tsoptions.GetParsedCommandLineOfConfigFile("/tsconfig.json", &core.CompilerOptions{}, nil, host, nil)
	assert.Equal(t, len(errors), 0, "Expected no errors in parsed command line")

	p := compiler.NewProgram(compiler.ProgramOptions{
		Config: parsed,
		Host:   host,
	})
	p.BindSourceFiles()
	c, done := p.GetTypeChecker(t.Context())
	defer done()
	file := p.GetSourceFile("/foo.ts")
	interfaceId := file.Statements.Nodes[0].Name()
	varId := file.Statements.Nodes[1].AsVariableStatement().DeclarationList.AsVariableDeclarationList().Declarations.Nodes[0].Name()
	propAccess := file.Statements.Nodes[2].Expression()
	nodes := []*ast.Node{interfaceId, varId, propAccess}
	for _, node := range nodes {
		symbol := c.GetSymbolAtLocation(node)
		if symbol == nil {
			t.Fatalf("Expected symbol to be non-nil")
		}
	}
}

func BenchmarkNewChecker(b *testing.B) {
	repo.SkipIfNoTypeScriptSubmodule(b)
	fs := osvfs.FS()
	fs = bundled.WrapFS(fs)

	rootPath := tspath.CombinePaths(tspath.NormalizeSlashes(repo.TypeScriptSubmodulePath()), "src", "compiler")

	host := compiler.NewCompilerHost(rootPath, fs, bundled.LibPath(), nil, nil)
	parsed, errors := tsoptions.GetParsedCommandLineOfConfigFile(tspath.CombinePaths(rootPath, "tsconfig.json"), &core.CompilerOptions{}, nil, host, nil)
	assert.Equal(b, len(errors), 0, "Expected no errors in parsed command line")
	p := compiler.NewProgram(compiler.ProgramOptions{
		Config: parsed,
		Host:   host,
	})

	b.ReportAllocs()

	for b.Loop() {
		checker.NewChecker(p, nil)
	}
}

func TestObjectNameCollisionInCommonJS(t *testing.T) {
	t.Parallel()

	makeProgram := func(module core.ModuleKind) *compiler.Program {
		fs := vfstest.FromMap(map[string]string{
			"/index.ts": `let Object = 0; export const x = 1;`,
		}, false /*useCaseSensitiveFileNames*/)
		fs = bundled.WrapFS(fs)

		opts := core.CompilerOptions{Module: module, Target: core.ScriptTargetESNext}
		return compiler.NewProgram(compiler.ProgramOptions{
			Config: &tsoptions.ParsedCommandLine{
				ParsedConfig: &core.ParsedOptions{
					FileNames:       []string{"/index.ts"},
					CompilerOptions: &opts,
				},
			},
			Host: compiler.NewCompilerHost("/", fs, bundled.LibPath(), nil, nil),
		})
	}

	hasReservedObjectDiagnostic := func(diags []*ast.Diagnostic) bool {
		for _, diag := range diags {
			if diag.Code() == diagnostics.Duplicate_identifier_0_Compiler_reserves_name_1_in_top_level_scope_of_a_module.Code() {
				return true
			}
		}
		return false
	}

	commonJSProgram := makeProgram(core.ModuleKindCommonJS)
	commonJSDiagnostics := commonJSProgram.GetSemanticDiagnostics(context.Background(), commonJSProgram.GetSourceFile("/index.ts"))
	assert.Assert(t, hasReservedObjectDiagnostic(commonJSDiagnostics), "expected CommonJS reserved-name diagnostic for Object")

	esNextProgram := makeProgram(core.ModuleKindESNext)
	esNextDiagnostics := esNextProgram.GetSemanticDiagnostics(context.Background(), esNextProgram.GetSourceFile("/index.ts"))
	assert.Assert(t, !hasReservedObjectDiagnostic(esNextDiagnostics), "did not expect CommonJS-only reserved-name diagnostic in ESNext module emit")
}
