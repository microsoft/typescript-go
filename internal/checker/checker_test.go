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

// TestKeyRemappingKeyofResult2 tests that index types for generic mapped types with name types
// don't crash (regression test for microsoft/TypeScript#56239)
func TestKeyRemappingKeyofResult2(t *testing.T) {
	t.Parallel()

	content := `// https://github.com/microsoft/TypeScript/issues/56239

type Values<T> = T[keyof T];

type ProvidedActor = {
  src: string;
  logic: unknown;
};

interface StateMachineConfig<TActors extends ProvidedActor> {
  invoke: {
    src: TActors["src"];
  };
}

declare function setup<TActors extends Record<string, unknown>>(_: {
  actors: {
    [K in keyof TActors]: TActors[K];
  };
}): {
  createMachine: (
    config: StateMachineConfig<
      Values<{
        [K in keyof TActors as K & string]: {
          src: K;
          logic: TActors[K];
        };
      }>
    >,
  ) => void;
};`

	fs := vfstest.FromMap(map[string]string{
		"/test.ts": content,
		"/tsconfig.json": `{
			"compilerOptions": {
				"strict": true,
				"noEmit": true
			},
			"files": ["test.ts"]
		}`,
	}, false)
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

	// The test passes if we can get a type checker without crashing
	assert.Assert(t, c != nil)
}

// TestMappedTypeAsClauseRecursiveNoCrash tests that recursive mapped types with as clauses
// don't crash when computing keyof (regression test for microsoft/TypeScript#60476)
func TestMappedTypeAsClauseRecursiveNoCrash(t *testing.T) {
	t.Parallel()

	content := `// https://github.com/microsoft/TypeScript/issues/60476

export type FlattenType<Source extends object, Target> = {
  [Key in keyof Source as Key extends string
    ? Source[Key] extends object
      ? ` + "`${Key}.${keyof FlattenType<Source[Key], Target> & string}`" + `
      : Key
    : never]-?: Target;
};

type FieldSelect = {
  table: string;
  field: string;
};

type Address = {
  postCode: string;
  description: string;
  address: string;
};

type User = {
  id: number;
  name: string;
  address: Address;
};

type FlattenedUser = FlattenType<User, FieldSelect>;
type FlattenedUserKeys = keyof FlattenType<User, FieldSelect>;

export type FlattenTypeKeys<Source extends object, Target> = keyof {
  [Key in keyof Source as Key extends string
    ? Source[Key] extends object
      ? ` + "`${Key}.${keyof FlattenType<Source[Key], Target> & string}`" + `
      : Key
    : never]-?: Target;
};

type FlattenedUserKeys2 = FlattenTypeKeys<User, FieldSelect>;`

	fs := vfstest.FromMap(map[string]string{
		"/test.ts": content,
		"/tsconfig.json": `{
			"compilerOptions": {
				"strict": true,
				"noEmit": true
			},
			"files": ["test.ts"]
		}`,
	}, false)
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

	// The test passes if we can get a type checker without crashing
	assert.Assert(t, c != nil)
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
		checker.NewChecker(p)
	}
}
