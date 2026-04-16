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
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func initDocumentLifecycleClient(t *testing.T, fs vfs.FS) *lsptestutil.LSPClient {
	t.Helper()

	onServerRequest := func(_ context.Context, req *lsproto.RequestMessage) *lsproto.ResponseMessage {
		switch req.Method {
		case lsproto.MethodWorkspaceConfiguration:
			return &lsproto.ResponseMessage{
				ID:      req.ID,
				JSONRPC: req.JSONRPC,
				Result:  []any{map[string]any{}},
			}
		case lsproto.MethodClientRegisterCapability, lsproto.MethodClientUnregisterCapability:
			return &lsproto.ResponseMessage{
				ID:      req.ID,
				JSONRPC: req.JSONRPC,
				Result:  lsproto.Null{},
			}
		default:
			return nil
		}
	}

	client, closeClient := lsptestutil.NewLSPClient(t, lsp.ServerOptions{
		Err:                io.Discard,
		Cwd:                "/home/projects",
		FS:                 bundled.WrapFS(fs),
		DefaultLibraryPath: bundled.LibPath(),
	}, onServerRequest)
	t.Cleanup(func() { _ = closeClient() })

	initMsg, _, ok := lsptestutil.SendRequest(t, client, lsproto.InitializeInfo, &lsproto.InitializeParams{})
	assert.Assert(t, ok, "expected initialize response")
	assert.Assert(t, initMsg.AsResponse().Error == nil)
	lsptestutil.SendNotification(t, client, lsproto.InitializedInfo, &lsproto.InitializedParams{})
	<-client.Server.InitComplete()

	return client
}

func TestCompletionAfterDuplicateDocumentClose(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	const fileName = "/home/projects/doc.ts"
	const text = "export const x = 1;\nx"

	fs := vfstest.FromMap(map[string]string{
		"/home/projects/tsconfig.json": `{"compilerOptions":{"module":"esnext","target":"esnext"}}`,
		fileName:                       text,
	}, false)
	client := initDocumentLifecycleClient(t, fs)
	uri := lsconv.FileNameToDocumentURI(fileName)

	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{
			Uri:        uri,
			LanguageId: "typescript",
			Version:    1,
			Text:       text,
		},
	})
	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidCloseInfo, &lsproto.DidCloseTextDocumentParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
	})

	msg, _, ok := lsptestutil.SendRequest(t, client, lsproto.TextDocumentCompletionInfo, &lsproto.CompletionParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
		Position:     lsproto.Position{Line: 1, Character: 1},
		Context:      &lsproto.CompletionContext{},
	})
	assert.Assert(t, ok, "expected response")
	assert.Assert(t, msg.AsResponse().Error == nil)

	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidCloseInfo, &lsproto.DidCloseTextDocumentParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
	})

	msg, _, ok = lsptestutil.SendRequest(t, client, lsproto.TextDocumentCompletionInfo, &lsproto.CompletionParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
		Position:     lsproto.Position{Line: 1, Character: 1},
		Context:      &lsproto.CompletionContext{},
	})
	assert.Assert(t, ok, "expected response")
	assert.Assert(t, msg.AsResponse().Error == nil)
}
