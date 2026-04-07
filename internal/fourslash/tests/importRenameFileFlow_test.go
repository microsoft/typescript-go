package fourslash_test

import (
	"slices"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
	"gotest.tools/v3/assert"
)

func fileRenameCapabilities() *lsproto.ClientCapabilities {
	capabilities := fourslash.GetDefaultCapabilities()
	capabilities.Workspace.WorkspaceEdit = &lsproto.WorkspaceEditClientCapabilities{
		DocumentChanges:    new(true),
		ResourceOperations: &[]lsproto.ResourceOperationKind{lsproto.ResourceOperationKindRename},
	}
	capabilities.Workspace.FileOperations = &lsproto.FileOperationClientCapabilities{
		WillRename: new(true),
	}
	return capabilities
}

func TestImportPathRenameReturnsRenameFileAndWillRenameEdits(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `// @Filename: /src/example.ts
import stuff from './[|stuff|].cts';
// @Filename: /src/stuff.cts
export = { name: "stuff" };
`

	f, done := fourslash.NewFourslash(t, fileRenameCapabilities(), content)
	defer done()
	f.Configure(t, &lsutil.UserPreferences{AllowRenameOfImportPath: core.TSTrue})
	f.GoToRangeStart(t, f.Ranges()[0])

	renameResult := f.RenameAtCaret(t, "renamed.cts")
	assert.Assert(t, renameResult.WorkspaceEdit != nil)
	assert.Assert(t, renameResult.WorkspaceEdit.DocumentChanges != nil)
	assert.Equal(t, len(*renameResult.WorkspaceEdit.DocumentChanges), 1)
	renameChange := (*renameResult.WorkspaceEdit.DocumentChanges)[0].RenameFile
	assert.Assert(t, renameChange != nil)
	assert.Equal(t, renameChange.OldUri, lsconv.FileNameToDocumentURI("/src/stuff.cts"))
	assert.Equal(t, renameChange.NewUri, lsconv.FileNameToDocumentURI("/src/renamed.cts"))

	willRenameResult := f.WillRenameFiles(t, &lsproto.FileRename{
		OldUri: string(lsconv.FileNameToDocumentURI("/src/stuff.cts")),
		NewUri: string(lsconv.FileNameToDocumentURI("/src/renamed.cts")),
	})
	assert.Assert(t, willRenameResult.WorkspaceEdit != nil)
	assert.Assert(t, willRenameResult.WorkspaceEdit.Changes != nil)

	edits := (*willRenameResult.WorkspaceEdit.Changes)[lsconv.FileNameToDocumentURI("/src/example.ts")]
	assert.Equal(t, len(edits), 1)
	assert.Equal(t, edits[0].NewText, "./renamed.cjs")
}

func TestImportPathDirectoryRenameReturnsRenameFileAndWillRenameEdits(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `// @Filename: /src/example.ts
import dir from './[|dir|]';
// @Filename: /src/dir/index.ts
export const x = 1;
`

	f, done := fourslash.NewFourslash(t, fileRenameCapabilities(), content)
	defer done()
	f.Configure(t, &lsutil.UserPreferences{AllowRenameOfImportPath: core.TSTrue})
	f.GoToRangeStart(t, f.Ranges()[0])

	renameResult := f.RenameAtCaret(t, "renamed")
	assert.Assert(t, renameResult.WorkspaceEdit != nil)
	assert.Assert(t, renameResult.WorkspaceEdit.DocumentChanges != nil)
	assert.Equal(t, len(*renameResult.WorkspaceEdit.DocumentChanges), 1)
	renameChange := (*renameResult.WorkspaceEdit.DocumentChanges)[0].RenameFile
	assert.Assert(t, renameChange != nil)
	assert.Equal(t, renameChange.OldUri, lsconv.FileNameToDocumentURI("/src/dir"))
	assert.Equal(t, renameChange.NewUri, lsconv.FileNameToDocumentURI("/src/renamed"))

	willRenameResult := f.WillRenameFiles(t, &lsproto.FileRename{
		OldUri: string(lsconv.FileNameToDocumentURI("/src/dir")),
		NewUri: string(lsconv.FileNameToDocumentURI("/src/renamed")),
	})
	assert.Assert(t, willRenameResult.WorkspaceEdit != nil)
	assert.Assert(t, willRenameResult.WorkspaceEdit.Changes != nil)

	edits := (*willRenameResult.WorkspaceEdit.Changes)[lsconv.FileNameToDocumentURI("/src/example.ts")]
	assert.Equal(t, len(edits), 1)
	assert.Equal(t, edits[0].NewText, "./renamed")
}

func TestWillRenameFilesUpdatesTsconfigAndTripleSlashReferences(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `// @Filename: /src/app.ts
/// <reference path="./old.ts" />
import { x } from "./old";
// @Filename: /src/old.ts
export const x = 1;
// @Filename: /tsconfig.json
{
  "files": ["src/app.ts", "src/old.ts"]
}
`

	f, done := fourslash.NewFourslash(t, fileRenameCapabilities(), content)
	defer done()

	willRenameResult := f.WillRenameFiles(t, &lsproto.FileRename{
		OldUri: string(lsconv.FileNameToDocumentURI("/src/old.ts")),
		NewUri: string(lsconv.FileNameToDocumentURI("/src/new.ts")),
	})
	assert.Assert(t, willRenameResult.WorkspaceEdit != nil)
	assert.Assert(t, willRenameResult.WorkspaceEdit.Changes != nil)

	appEdits := (*willRenameResult.WorkspaceEdit.Changes)[lsconv.FileNameToDocumentURI("/src/app.ts")]
	assert.Equal(t, len(appEdits), 2)
	newTexts := []string{appEdits[0].NewText, appEdits[1].NewText}
	slices.Sort(newTexts)
	assert.DeepEqual(t, newTexts, []string{"./new", "./new.ts"})

	tsconfigEdits := (*willRenameResult.WorkspaceEdit.Changes)[lsconv.FileNameToDocumentURI("/tsconfig.json")]
	assert.Equal(t, len(tsconfigEdits), 1)
	assert.Equal(t, tsconfigEdits[0].NewText, "src/new.ts")
}

func TestWillRenameFilesUpdatesProjectReferenceConsumer(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `// @Filename: /solution/a/old.ts
export const x = 1;
// @Filename: /solution/b/app.ts
import { x } from "../a/old";
// @Filename: /solution/a/tsconfig.json
{
  "compilerOptions": {
    "composite": true
  },
  "files": ["old.ts"]
}
// @Filename: /solution/b/tsconfig.json
{
  "references": [{ "path": "../a" }],
  "files": ["app.ts"]
}
`

	f, done := fourslash.NewFourslash(t, fileRenameCapabilities(), content)
	defer done()

	willRenameResult := f.WillRenameFiles(t, &lsproto.FileRename{
		OldUri: string(lsconv.FileNameToDocumentURI("/solution/a/old.ts")),
		NewUri: string(lsconv.FileNameToDocumentURI("/solution/a/new.ts")),
	})
	assert.Assert(t, willRenameResult.WorkspaceEdit != nil)
	assert.Assert(t, willRenameResult.WorkspaceEdit.Changes != nil)

	appEdits := (*willRenameResult.WorkspaceEdit.Changes)[lsconv.FileNameToDocumentURI("/solution/b/app.ts")]
	assert.Equal(t, len(appEdits), 1)
	assert.Equal(t, appEdits[0].NewText, "../a/new")
}

func TestWillRenameFilesUpdatesSiblingProjectLoadedViaSolutionRoot(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `// @Filename: /solution/a/old.ts
export const x = 1;
// @Filename: /solution/b/app.ts
import { x } from "../a/old";
// @Filename: /solution/tsconfig.json
{
  "files": [],
  "references": [
    { "path": "./a" },
    { "path": "./b" }
  ]
}
// @Filename: /solution/a/tsconfig.json
{
  "compilerOptions": {
    "composite": true
  },
  "files": ["old.ts"]
}
// @Filename: /solution/b/tsconfig.json
{
  "files": ["app.ts"]
}
`

	f, done := fourslash.NewFourslash(t, fileRenameCapabilities(), content)
	defer done()

	willRenameResult := f.WillRenameFiles(t, &lsproto.FileRename{
		OldUri: string(lsconv.FileNameToDocumentURI("/solution/a/old.ts")),
		NewUri: string(lsconv.FileNameToDocumentURI("/solution/a/new.ts")),
	})
	assert.Assert(t, willRenameResult.WorkspaceEdit != nil)
	assert.Assert(t, willRenameResult.WorkspaceEdit.Changes != nil)

	appEdits := (*willRenameResult.WorkspaceEdit.Changes)[lsconv.FileNameToDocumentURI("/solution/b/app.ts")]
	assert.Equal(t, len(appEdits), 1)
	assert.Equal(t, appEdits[0].NewText, "../a/new")
}

func TestWillRenameFilesUpdatesSiblingProjectWhenUnrelatedFileIsInitiallyActive(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `// @Filename: /solution/c/unrelated.ts
export const active = 1;
// @Filename: /solution/a/old.ts
export const x = 1;
// @Filename: /solution/b/app.ts
import { x } from "../a/old";
// @Filename: /solution/tsconfig.json
{
  "files": [],
  "references": [
    { "path": "./a" },
    { "path": "./b" },
    { "path": "./c" }
  ]
}
// @Filename: /solution/a/tsconfig.json
{
  "compilerOptions": {
    "composite": true
  },
  "files": ["old.ts"]
}
// @Filename: /solution/b/tsconfig.json
{
  "files": ["app.ts"]
}
// @Filename: /solution/c/tsconfig.json
{
  "files": ["unrelated.ts"]
}
`

	f, done := fourslash.NewFourslash(t, fileRenameCapabilities(), content)
	defer done()

	willRenameResult := f.WillRenameFiles(t, &lsproto.FileRename{
		OldUri: string(lsconv.FileNameToDocumentURI("/solution/a/old.ts")),
		NewUri: string(lsconv.FileNameToDocumentURI("/solution/a/new.ts")),
	})
	assert.Assert(t, willRenameResult.WorkspaceEdit != nil)
	assert.Assert(t, willRenameResult.WorkspaceEdit.Changes != nil)

	appEdits := (*willRenameResult.WorkspaceEdit.Changes)[lsconv.FileNameToDocumentURI("/solution/b/app.ts")]
	assert.Equal(t, len(appEdits), 1)
	assert.Equal(t, appEdits[0].NewText, "../a/new")
}

func TestWillRenameFilesUpdatesSymlinkedPackageConsumer(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `// @Filename: /packages/project-b/old.ts
export const x = 1;
// @Filename: /packages/project-b/package.json
{
  "name": "project-b",
  "version": "1.0.0"
}
// @Filename: /packages/project-b/tsconfig.json
{
  "compilerOptions": {
    "composite": true,
    "module": "commonjs"
  },
  "files": ["old.ts"]
}
// @Filename: /packages/project-a/app.ts
import { x } from "project-b/old";
// @Filename: /packages/project-a/package.json
{
  "name": "project-a",
  "dependencies": {
    "project-b": "*"
  }
}
// @Filename: /packages/project-a/tsconfig.json
{
  "compilerOptions": {
    "module": "commonjs"
  },
  "files": ["app.ts"]
}
// @link: /packages/project-b -> /packages/project-a/node_modules/project-b`

	f, done := fourslash.NewFourslash(t, fileRenameCapabilities(), content)
	defer done()

	willRenameResult := f.WillRenameFiles(t, &lsproto.FileRename{
		OldUri: string(lsconv.FileNameToDocumentURI("/packages/project-b/old.ts")),
		NewUri: string(lsconv.FileNameToDocumentURI("/packages/project-b/new.ts")),
	})
	assert.Assert(t, willRenameResult.WorkspaceEdit != nil)
	assert.Assert(t, willRenameResult.WorkspaceEdit.Changes != nil)

	appEdits := (*willRenameResult.WorkspaceEdit.Changes)[lsconv.FileNameToDocumentURI("/packages/project-a/app.ts")]
	assert.Equal(t, len(appEdits), 1)
	assert.Equal(t, appEdits[0].NewText, "project-b/new")
}

func TestImportTypePathRenameReturnsRenameFile(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `// @module: commonjs
// @Filename: /a.ts
export = 0;
// @Filename: /b.ts
const x: import("[|./a|]") = 0;
`

	f, done := fourslash.NewFourslash(t, fileRenameCapabilities(), content)
	defer done()

	prefsTrue := &lsutil.UserPreferences{AllowRenameOfImportPath: core.TSTrue}
	prefsFalse := &lsutil.UserPreferences{AllowRenameOfImportPath: core.TSFalse}

	f.Configure(t, prefsTrue)
	f.GoToRangeStart(t, f.Ranges()[0])

	renameResult := f.RenameAtCaret(t, "renamed.ts")
	assert.Assert(t, renameResult.WorkspaceEdit != nil)
	assert.Assert(t, renameResult.WorkspaceEdit.DocumentChanges != nil)
	assert.Equal(t, len(*renameResult.WorkspaceEdit.DocumentChanges), 1)
	renameChange := (*renameResult.WorkspaceEdit.DocumentChanges)[0].RenameFile
	assert.Assert(t, renameChange != nil)
	assert.Equal(t, renameChange.OldUri, lsconv.FileNameToDocumentURI("/a.ts"))
	assert.Equal(t, renameChange.NewUri, lsconv.FileNameToDocumentURI("/renamed.ts"))

	f.Configure(t, prefsFalse)
	f.GoToRangeStart(t, f.Ranges()[0])
	f.VerifyRenameFailed(t, prefsFalse)
}

func TestGlobalImportRenameStillFails(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `// @allowJs: true
// @module: commonjs
// @Filename: /node_modules/global/index.d.ts
export const x: number;
// @Filename: /c.js
const global = require("/*global*/global");
`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	prefsTrue := &lsutil.UserPreferences{
		IncludeCompletionsForModuleExports:    core.TSTrue,
		IncludeCompletionsForImportStatements: core.TSTrue,
		AllowRenameOfImportPath:               core.TSTrue,
	}

	f.Configure(t, prefsTrue)
	f.GoToMarker(t, "global")
	f.VerifyRenameFailed(t, prefsTrue)
}

func TestImportPathRenameFailsWithoutFileRenameClientSupport(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `// @Filename: /src/example.ts
import stuff from './[|stuff|].cts';
// @Filename: /src/stuff.cts
export = { name: "stuff" };
`

	capabilities := fileRenameCapabilities()
	capabilities.Workspace.FileOperations = nil

	f, done := fourslash.NewFourslash(t, capabilities, content)
	defer done()
	f.Configure(t, &lsutil.UserPreferences{AllowRenameOfImportPath: core.TSTrue})
	f.GoToRangeStart(t, f.Ranges()[0])
	f.VerifyRenameFailed(t, &lsutil.UserPreferences{AllowRenameOfImportPath: core.TSTrue})
}
