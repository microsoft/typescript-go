package checker_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/testutil"
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

func TestImportHelpersAfterProgramUpdateWithoutSyntheticImportSpecifier(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on import helpers after program update")

	const config = `{
		"compilerOptions": {
			"target": "es2015",
			"module": "commonjs",
			"experimentalDecorators": true,
			"importHelpers": true
		},
		"files": ["foo.ts"]
	}`
	const original = `declare function dec(value: Function): void;
class C {}`
	const updated = `declare function dec(value: Function): void;
@dec
export class C {}`
	const tslib = `export declare function __decorate(...args: any[]): any;`

	fs := bundled.WrapFS(vfstest.FromMap(map[string]string{
		"/foo.ts":                          original,
		"/node_modules/tslib/package.json": `{"name":"tslib","typings":"tslib.d.ts"}`,
		"/node_modules/tslib/tslib.d.ts":   tslib,
		"/node_modules/tslib/tslib.js":     "",
		"/tsconfig.json":                   config,
	}, false /*useCaseSensitiveFileNames*/))
	host := compiler.NewCompilerHost("/", fs, bundled.LibPath(), nil, nil)
	parsed, errors := tsoptions.GetParsedCommandLineOfConfigFile("/tsconfig.json", &core.CompilerOptions{}, nil, host, nil)
	assert.Equal(t, len(errors), 0, "Expected no errors in parsed command line")

	p := compiler.NewProgram(compiler.ProgramOptions{
		Config: parsed,
		Host:   host,
	})
	path := p.GetSourceFile("/foo.ts").Path()

	updatedFS := bundled.WrapFS(vfstest.FromMap(map[string]string{
		"/foo.ts":                          updated,
		"/node_modules/tslib/package.json": `{"name":"tslib","typings":"tslib.d.ts"}`,
		"/node_modules/tslib/tslib.d.ts":   tslib,
		"/node_modules/tslib/tslib.js":     "",
		"/tsconfig.json":                   config,
	}, false /*useCaseSensitiveFileNames*/))
	updatedHost := compiler.NewCompilerHost("/", updatedFS, bundled.LibPath(), nil, nil)
	updatedProgram, _ := p.UpdateProgram(path, updatedHost, nil)
	updatedFile := updatedProgram.GetSourceFile("/foo.ts")
	diagnostics := updatedProgram.GetSemanticDiagnostics(t.Context(), updatedFile)
	assert.Equal(t, len(diagnostics), 0, "Expected no semantic diagnostics")
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
