package lsp_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
	"gotest.tools/v3/assert"
)

func TestWillRenameFilesAfterWatchedFileRenameInCompositeProject(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	client, fs := initMutableLSPClient(t, map[string]string{
		"/root/package.json": `{
			"name": "tsgo-test",
			"version": "0.1.0",
			"type": "module"
		}`,
		"/root/tsconfig.base.json": `{
			"compilerOptions": {
				"composite": true,
				"module": "nodenext",
				"target": "esnext",
				"types": [],
				"sourceMap": true,
				"declaration": true,
				"declarationMap": true,
				"strict": true,
				"jsx": "react-jsx",
				"verbatimModuleSyntax": true,
				"isolatedModules": true,
				"noUncheckedSideEffectImports": true,
				"moduleDetection": "force",
				"skipLibCheck": true
			}
		}`,
		"/root/tsconfig.json": `{
			"files": [],
			"references": [
				{ "path": "./projects/a/" }
			]
		}`,
		"/root/projects/a/tsconfig.json": `{
			"extends": "../../tsconfig.base.json"
		}`,
		"/root/projects/a/b.ts": "export const x = 0\n",
	}, &lsutil.UserPreferences{})

	oldURI := lsconv.FileNameToDocumentURI("/root/projects/a/b.ts")
	newURI := lsconv.FileNameToDocumentURI("/root/projects/a/b1.ts")

	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{Uri: oldURI, LanguageId: "typescript", Text: "export const x = 0\n"},
	})

	assert.NilError(t, fs.Remove("root/projects/a/b.ts"))
	assert.NilError(t, fs.WriteFile("root/projects/a/b1.ts", "export const x = 0\n", 0o666))
	lsptestutil.SendNotification(t, client, lsproto.WorkspaceDidChangeWatchedFilesInfo, &lsproto.DidChangeWatchedFilesParams{
		Changes: []*lsproto.FileEvent{
			{Uri: oldURI, Type: lsproto.FileChangeTypeDeleted},
			{Uri: newURI, Type: lsproto.FileChangeTypeCreated},
		},
	})

	msg, _, ok := lsptestutil.SendRequest(t, client, lsproto.WorkspaceWillRenameFilesInfo, &lsproto.RenameFilesParams{
		Files: []*lsproto.FileRename{{
			OldUri: string(oldURI),
			NewUri: string(newURI),
		}},
	})
	assert.Assert(t, ok, "expected response")
	assert.Assert(t, msg.AsResponse().Error == nil)
}
