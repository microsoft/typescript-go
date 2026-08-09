package checker_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
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

func TestGetTypeAtLocationTypeOnlyImportClause(t *testing.T) {
	t.Parallel()

	// A type-only import clause is a type declaration, but it only has a symbol
	// when it declares a default import name. Every other form leaves the clause
	// symbol-less, and getTypeOfNode used to hand that nil symbol straight to
	// getDeclaredTypeOfSymbol.
	testCases := []struct {
		name    string
		content string
	}{
		{"named", `import type { A } from "./x";`},
		{"namespace", `import type * as ns from "./x";`},
		{"empty", `import type {} from "./x";`},
		{"default", `import type A from "./x";`},
		{"defaultAndNamed", `import type A, { B } from "./x";`},
		{"unresolvedModule", `import type { A } from "./nonexistent";`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			fs := vfstest.FromMap(map[string]string{
				"/foo.ts": testCase.content,
				"/x.ts":   "export interface A { a: string }\nexport interface B { b: string }",
				"/tsconfig.json": `
						{
							"compilerOptions": {},
							"files": ["foo.ts", "x.ts"]
						}
					`,
			}, false /*useCaseSensitiveFileNames*/)
			fs = bundled.WrapFS(fs)

			host := compiler.NewCompilerHost("/", fs, bundled.LibPath(), nil, nil)

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
			importClause := file.Statements.Nodes[0].AsImportDeclaration().ImportClause
			if importClause == nil {
				t.Fatalf("Expected an import clause")
			}

			// Should return a type rather than panicking.
			if c.GetTypeAtLocation(importClause) == nil {
				t.Fatalf("Expected GetTypeAtLocation to return a type")
			}
		})
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
