package autoimport

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
	"github.com/microsoft/typescript-go/internal/packagejson"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

type fakeCloneHost struct {
	fs          vfs.FS
	sourceFiles map[tspath.Path]*ast.SourceFile
}

func (h *fakeCloneHost) FS() vfs.FS                  { return h.fs }
func (h *fakeCloneHost) GetCurrentDirectory() string { return "/" }
func (h *fakeCloneHost) GetDefaultProject(path tspath.Path) (tspath.Path, *compiler.Program) {
	return "", nil
}

func (h *fakeCloneHost) GetProgramForProject(projectPath tspath.Path) *compiler.Program { return nil }

func (h *fakeCloneHost) GetPackageJson(fileName string) *packagejson.InfoCacheEntry { return nil }

func (h *fakeCloneHost) GetSourceFile(fileName string, path tspath.Path) *ast.SourceFile {
	return h.sourceFiles[path]
}
func (h *fakeCloneHost) Dispose() {}

var _ RegistryCloneHost = (*fakeCloneHost)(nil)

// Regression test for microsoft/typescript-go#4322.
//
// During auto-import export extraction, the checker is built on top of an
// aliasResolver standing in for a real program. This file has a type error, and
// extracting exports should still complete without crashing.
func TestAliasResolverGetDiagnosticsDoesNotPanic(t *testing.T) {
	t.Parallel()

	const fileName = "/pkg/index.ts"
	text := "declare function f(arg: { a: string }): () => void;\nexport const x = f({ a: 1 });\n"

	fs := vfstest.FromMap(map[string]string{fileName: text}, true /*useCaseSensitiveFileNames*/)
	host := &fakeCloneHost{fs: fs}

	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: fileName,
		Path:     tspath.Path(fileName),
	}, text, core.ScriptKindTS)
	binder.BindSourceFile(sourceFile)

	resolver := module.NewResolver(host, core.EmptyCompilerOptions, "", "")
	r := newAliasResolver(
		[]*ast.SourceFile{sourceFile},
		nil,
		host,
		resolver,
		func(f string) tspath.Path { return tspath.Path(f) },
		func(ast.HasFileName, string) {},
	)

	ch, _ := checker.NewChecker(r, nil)

	// Type-checking this file's diagnostics must not panic.
	ch.GetDiagnostics(context.Background(), sourceFile)
}

// Regression test for microsoft/typescript-go#4481.
//
// When a file exports a value whose type references a symbol from another module,
// the checker's node builder calls GetSourceOfProjectReferenceIfOutputIncluded
// to generate module specifiers. The aliasResolver must not panic.
func TestAliasResolverGetSourceOfProjectReferenceDoesNotPanic(t *testing.T) {
	t.Parallel()

	files := map[string]string{
		"/pkg/types.ts":    "export interface Foo { x: number; }\n",
		"/pkg/index.ts":    "import { Foo } from './types';\nexport function getFoo(): Foo { return { x: 1 }; }\n",
		"/pkg/consumer.ts": "",
	}

	fs := vfstest.FromMap(files, true /*useCaseSensitiveFileNames*/)
	sourceFiles := make(map[tspath.Path]*ast.SourceFile)
	for fileName, text := range files {
		path := tspath.Path(fileName)
		sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
			FileName: fileName,
			Path:     path,
		}, text, core.ScriptKindTS)
		binder.BindSourceFile(sourceFile)
		sourceFiles[path] = sourceFile
	}
	host := &fakeCloneHost{fs: fs, sourceFiles: sourceFiles}
	indexFile := sourceFiles["/pkg/index.ts"]
	consumerFile := sourceFiles["/pkg/consumer.ts"]

	resolver := module.NewResolver(host, core.EmptyCompilerOptions, "/pkg", "")
	r := newAliasResolver(
		[]*ast.SourceFile{sourceFiles["/pkg/types.ts"], indexFile, consumerFile},
		nil,
		host,
		resolver,
		func(f string) tspath.Path { return tspath.Path(f) },
		func(ast.HasFileName, string) {},
	)

	ch, _ := checker.NewChecker(r, nil)

	getFooSymbol := indexFile.Symbol.Exports["getFoo"]
	tpe := ch.GetTypeOfSymbol(getFooSymbol)

	// Serializing getFoo's type in a different file has to synthesize a module
	// specifier for Foo, which calls GetSourceOfProjectReferenceIfOutputIncluded.
	// This must not panic.
	ch.TypeToTypeNode(tpe, consumerFile.AsNode(), nodebuilder.FlagsNone, nil)
}
