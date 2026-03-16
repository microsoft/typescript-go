package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
	"gotest.tools/v3/assert"
)

func TestGetEditsForFileRenameUpdatesImports(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	f, done := fourslash.NewFourslash(t, fileRenameCapabilities(), `
// @Filename: /tsconfig.json
{"compilerOptions":{"module":"esnext","moduleResolution":"bundler","target":"esnext","noLib":true},"include":["**/*.ts"]}
// @Filename: /shared/util.ts
export const util = 1;
// @Filename: /foo.ts
import { util } from './shared/util';
export const foo = util;
// @Filename: /main.ts
import { foo } from "./foo";
foo;
`)
	defer done()

	result := f.GetEditsForFileRename(t, "/foo.ts", "/nested/foo.ts")
	changes := requireWorkspaceChanges(t, result)

	assert.Equal(t, len(changes), 2)
	assert.Equal(t, f.ApplyTextEdits(t, "/main.ts", changes[lsconv.FileNameToDocumentURI("/main.ts")]), "import { foo } from \"./nested/foo\";\nfoo;\n")
	assert.Equal(t, f.ApplyTextEdits(t, "/foo.ts", changes[lsconv.FileNameToDocumentURI("/foo.ts")]), "import { util } from '../shared/util';\nexport const foo = util;\n")
}

func TestGetEditsForFileRenameUnaffectedNonRelativePath(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	f, done := fourslash.NewFourslash(t, fileRenameCapabilities(), `
// @Filename: /sub/a.ts
export const a = 1;
// @Filename: /sub/b.ts
import { a } from "sub/a";
// @Filename: /tsconfig.json
{
    "compilerOptions": {
        "baseUrl": "."
    }
}
`)
	defer done()

	result := f.GetEditsForFileRename(t, "/sub/b.ts", "/sub/c/d.ts")
	assert.Assert(t, result.WorkspaceEdit == nil)
}

func TestGetEditsForFileRenameHandlesBatchRenames(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	f, done := fourslash.NewFourslash(t, fileRenameCapabilities(), `
// @Filename: /tsconfig.json
{"compilerOptions":{"module":"esnext","moduleResolution":"bundler","target":"esnext","noLib":true},"include":["**/*.ts"]}
// @Filename: /a.ts
export const a = 1;
// @Filename: /b.ts
export const b = 1;
// @Filename: /main.ts
import { a } from "./a";
import { b } from './b';
a;
b;
`)
	defer done()

	result := f.GetEditsForFileRenames(t, []*lsproto.FileRename{
		{
			OldUri: string(lsconv.FileNameToDocumentURI("/a.ts")),
			NewUri: string(lsconv.FileNameToDocumentURI("/nested/a.ts")),
		},
		{
			OldUri: string(lsconv.FileNameToDocumentURI("/b.ts")),
			NewUri: string(lsconv.FileNameToDocumentURI("/nested/b.ts")),
		},
	})
	changes := requireWorkspaceChanges(t, result)

	assert.Equal(t, len(changes), 1)
	assert.Equal(t, f.ApplyTextEdits(t, "/main.ts", changes[lsconv.FileNameToDocumentURI("/main.ts")]), "import { a } from \"./nested/a\";\nimport { b } from './nested/b';\na;\nb;\n")
}

func fileRenameCapabilities() *lsproto.ClientCapabilities {
	capabilities := fourslash.GetDefaultCapabilities()
	capabilities.Workspace.FileOperations = &lsproto.FileOperationClientCapabilities{
		DidRename:  ptr(true),
		WillRename: ptr(true),
	}
	return capabilities
}

func requireWorkspaceChanges(t *testing.T, result lsproto.WorkspaceEditOrNull) map[lsproto.DocumentUri][]*lsproto.TextEdit {
	t.Helper()
	if result.WorkspaceEdit == nil || result.WorkspaceEdit.Changes == nil {
		t.Fatal("expected workspace edits")
	}
	return *result.WorkspaceEdit.Changes
}

func ptr[T any](v T) *T {
	return &v
}
