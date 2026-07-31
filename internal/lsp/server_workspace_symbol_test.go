package lsp_test

import (
	"context"
	"io"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func TestWorkspaceSymbolsCurrentProject(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	files := map[string]string{
		"/home/projects/a/tsconfig.json": `{}`,
		"/home/projects/a/index.ts":      `export function fromA() {}`,
		"/home/projects/b/tsconfig.json": `{}`,
		"/home/projects/b/index.ts":      `export function fromB() {}`,
	}
	fs := bundled.WrapFS(vfstest.FromMap(files, false))
	onServerRequest := func(_ context.Context, req *lsproto.RequestMessage) *lsproto.ResponseMessage {
		if req.Method == lsproto.MethodClientRegisterCapability || req.Method == lsproto.MethodClientUnregisterCapability {
			return &lsproto.ResponseMessage{ID: req.ID, JSONRPC: req.JSONRPC, Result: lsproto.Null{}}
		}
		return nil
	}
	client, closeClient := lsptestutil.NewLSPClient(t, lsp.ServerOptions{
		Err:                io.Discard,
		Cwd:                "/home/projects",
		FS:                 fs,
		DefaultLibraryPath: bundled.LibPath(),
	}, onServerRequest)
	t.Cleanup(func() { _ = closeClient() })

	initMsg, _, ok := lsptestutil.SendRequest(t, client, lsproto.InitializeInfo, &lsproto.InitializeParams{
		Capabilities: &lsproto.ClientCapabilities{},
	})
	assert.Assert(t, ok && initMsg.AsResponse().Error == nil, "Initialize failed")
	lsptestutil.SendNotification(t, client, lsproto.InitializedInfo, &lsproto.InitializedParams{})
	<-client.Server.InitComplete()

	for _, file := range []string{"/home/projects/a/index.ts", "/home/projects/b/index.ts"} {
		lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
			TextDocument: &lsproto.TextDocumentItem{
				Uri:        lsconv.FileNameToDocumentURI(file),
				LanguageId: lsproto.LanguageKindTypeScript,
				Version:    1,
				Text:       files[file],
			},
		})
	}

	workspaceSymbolParams := &lsproto.WorkspaceSymbolParams{
		Query:        "from",
		TextDocument: &lsproto.TextDocumentIdentifier{Uri: lsconv.FileNameToDocumentURI("/home/projects/a/index.ts")},
	}
	_, allProjects, ok := lsptestutil.SendRequest(t, client, lsproto.WorkspaceSymbolInfo, workspaceSymbolParams)
	assert.Assert(t, ok)
	assert.Equal(t, len(*allProjects.SymbolInformations), 2)

	lsptestutil.SendNotification(t, client, lsproto.WorkspaceDidChangeConfigurationInfo, &lsproto.DidChangeConfigurationParams{
		Settings: map[string]any{
			"js/ts": map[string]any{
				"workspaceSymbols": map[string]any{
					"scope": "currentProject",
				},
			},
		},
	})

	_, currentProject, ok := lsptestutil.SendRequest(t, client, lsproto.WorkspaceSymbolInfo, workspaceSymbolParams)
	assert.Assert(t, ok)
	assert.Equal(t, len(*currentProject.SymbolInformations), 1)
	assert.Equal(t, (*currentProject.SymbolInformations)[0].Name, "fromA")
}
