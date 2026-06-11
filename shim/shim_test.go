package shim_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/microsoft/typescript-go/shim/ast"
	"github.com/microsoft/typescript-go/shim/bundled"
	"github.com/microsoft/typescript-go/shim/compiler"
	"github.com/microsoft/typescript-go/shim/core"
	"github.com/microsoft/typescript-go/shim/tsoptions"
	"github.com/microsoft/typescript-go/shim/vfs"
	"github.com/microsoft/typescript-go/shim/vfs/osvfs"
)

// TestShimSmoke exercises the full external consumer flow through the shim
// packages only: parse a tsconfig, build a Program against the bundled libs,
// obtain a type checker, and walk the AST. It is the canary that fails an
// upstream-sync PR when a re-exported API moves.
func TestShimSmoke(t *testing.T) {
	t.Parallel()

	dir := filepath.ToSlash(t.TempDir())
	writeFile(t, filepath.Join(dir, "tsconfig.json"), `{
		"compilerOptions": {
			"strict": true,
			"target": "es2022",
			"module": "esnext",
			"noEmit": true
		},
		"include": ["**/*.schema.ts"]
	}`)
	writeFile(t, filepath.Join(dir, "person.schema.ts"), `
export class Person {
	name: string = "";
	greet(): string {
		return "hello " + this.name;
	}
}
`)

	var fs vfs.FS = bundled.WrapFS(osvfs.FS())
	host := compiler.NewCachedFSCompilerHost(dir, fs, bundled.LibPath(), nil, nil)

	configFileName := dir + "/tsconfig.json"
	config, diags := tsoptions.GetParsedCommandLineOfConfigFile(configFileName, nil, nil, host, nil)
	if len(diags) != 0 {
		t.Fatalf("unexpected tsconfig diagnostics: %v", diagnosticMessages(diags))
	}
	if config == nil {
		t.Fatal("expected a parsed command line")
	}

	program := compiler.NewProgram(compiler.ProgramOptions{
		Host:           host,
		Config:         config,
		SingleThreaded: core.TSTrue,
	})

	ctx := context.Background()
	checker, done := program.GetTypeChecker(ctx)
	defer done()
	if checker == nil {
		t.Fatal("expected a type checker")
	}

	var schemaFile *ast.SourceFile
	for _, file := range program.SourceFiles() {
		if filepath.Base(file.FileName()) == "person.schema.ts" {
			schemaFile = file
		}
	}
	if schemaFile == nil {
		t.Fatal("expected person.schema.ts in the program")
	}

	if semDiags := program.GetSemanticDiagnostics(ctx, schemaFile); len(semDiags) != 0 {
		t.Fatalf("unexpected semantic diagnostics: %v", diagnosticMessages(semDiags))
	}

	var classNode *ast.Node
	for _, statement := range schemaFile.Statements.Nodes {
		if statement.Kind == ast.KindClassDeclaration {
			classNode = statement
		}
	}
	if classNode == nil {
		t.Fatal("expected a class declaration in person.schema.ts")
	}
	if name := classNode.Name(); name == nil || name.Text() != "Person" {
		t.Fatalf("expected class named Person, got %v", classNode.Name())
	}
}

func writeFile(t *testing.T, path string, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func diagnosticMessages(diags []*ast.Diagnostic) []string {
	messages := make([]string, 0, len(diags))
	for _, diag := range diags {
		messages = append(messages, diag.String())
	}
	return messages
}
