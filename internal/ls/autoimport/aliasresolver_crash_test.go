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
	"github.com/microsoft/typescript-go/internal/packagejson"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

type fakeCloneHost struct {
	fs vfs.FS
}

func (h *fakeCloneHost) FS() vfs.FS                  { return h.fs }
func (h *fakeCloneHost) GetCurrentDirectory() string { return "/" }
func (h *fakeCloneHost) GetDefaultProject(path tspath.Path) (tspath.Path, *compiler.Program) {
	return "", nil
}

func (h *fakeCloneHost) GetProgramForProject(projectPath tspath.Path) *compiler.Program { return nil }

func (h *fakeCloneHost) GetPackageJson(fileName string) *packagejson.InfoCacheEntry { return nil }

func (h *fakeCloneHost) GetSourceFile(fileName string, path tspath.Path) *ast.SourceFile { return nil }
func (h *fakeCloneHost) Dispose()                                                        {}

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
		"/pkg/types.ts": "export interface Foo { x: number; }\n",
		"/pkg/index.ts": "import { Foo } from './types';\nexport function getFoo(): Foo { return { x: 1 }; }\n",
	}

	fs := vfstest.FromMap(files, true /*useCaseSensitiveFileNames*/)
	host := &fakeCloneHost{fs: fs}

	typesFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/pkg/types.ts",
		Path:     "/pkg/types.ts",
	}, files["/pkg/types.ts"], core.ScriptKindTS)
	binder.BindSourceFile(typesFile)

	indexFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/pkg/index.ts",
		Path:     "/pkg/index.ts",
	}, files["/pkg/index.ts"], core.ScriptKindTS)
	binder.BindSourceFile(indexFile)

	resolver := module.NewResolver(host, core.EmptyCompilerOptions, "/pkg", "")
	r := newAliasResolver(
		[]*ast.SourceFile{typesFile, indexFile},
		nil,
		host,
		resolver,
		func(f string) tspath.Path { return tspath.Path(f) },
		func(ast.HasFileName, string) {},
	)

	ch, _ := checker.NewChecker(r, nil)

	// Getting diagnostics triggers the node builder which calls
	// GetSourceOfProjectReferenceIfOutputIncluded for module specifier generation.
	// This must not panic.
	ch.GetDiagnostics(context.Background(), indexFile)
}
