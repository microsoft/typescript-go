package lsp_test

import (
	"context"
	"io"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func initProjectInfoClient(t *testing.T, files map[string]string) *lsptestutil.LSPClient {
	t.Helper()

	fs := bundled.WrapFS(vfstest.FromMap(files, false))

	onServerRequest := func(_ context.Context, req *lsproto.RequestMessage) *lsproto.ResponseMessage {
		switch req.Method {
		case lsproto.MethodClientRegisterCapability, lsproto.MethodClientUnregisterCapability, lsproto.MethodWindowWorkDoneProgressCreate:
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

	return client
}

func TestProjectInfoConfiguredProject(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	client := initProjectInfoClient(t, map[string]string{
		"/home/projects/tsconfig.json": `{}`,
		"/home/projects/index.ts":      "export const x = 1;",
	})

	uri := lsproto.DocumentUri("file:///home/projects/index.ts")
	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{Uri: uri, LanguageId: "typescript", Text: "export const x = 1;"},
	})

	msg, resp, ok := lsptestutil.SendRequest(t, client, lsproto.CustomProjectInfoInfo, &lsproto.ProjectInfoParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
	})
	assert.Assert(t, ok, "expected a response")
	assert.Assert(t, msg.AsResponse().Error == nil)
	assert.Equal(t, resp.ConfigFilePath, "/home/projects/tsconfig.json")
}

func TestProjectInfoInferredProject(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	client := initProjectInfoClient(t, map[string]string{
		"/home/projects/index.ts": "export const x = 1;",
	})

	uri := lsproto.DocumentUri("file:///home/projects/index.ts")
	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{Uri: uri, LanguageId: "typescript", Text: "export const x = 1;"},
	})

	msg, resp, ok := lsptestutil.SendRequest(t, client, lsproto.CustomProjectInfoInfo, &lsproto.ProjectInfoParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
	})
	assert.Assert(t, ok, "expected a response")
	assert.Assert(t, msg.AsResponse().Error == nil)
	assert.Equal(t, resp.ConfigFilePath, "")
}

func TestProjectInfoUntitledFile(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	// Test that untitled files get an inferred project and don't produce errors.
	client := initProjectInfoClient(t, map[string]string{})

	uri := lsproto.DocumentUri("untitled:Untitled-1")
	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{Uri: uri, LanguageId: "typescript", Text: ""},
	})

	// custom/projectInfo should succeed without error for an untitled file.
	msg, resp, ok := lsptestutil.SendRequest(t, client, lsproto.CustomProjectInfoInfo, &lsproto.ProjectInfoParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
	})
	assert.Assert(t, ok, "expected a response")
	assert.Assert(t, msg.AsResponse().Error == nil, "expected no error for untitled file, got: %v", msg.AsResponse().Error)
	assert.Equal(t, resp.ConfigFilePath, "")

	// textDocument/diagnostic should also succeed without error.
	diagMsg, _, diagOk := lsptestutil.SendRequest(t, client, lsproto.TextDocumentDiagnosticInfo, &lsproto.DocumentDiagnosticParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
	})
	assert.Assert(t, diagOk, "expected a diagnostic response")
	assert.Assert(t, diagMsg.AsResponse().Error == nil, "expected no error for untitled file diagnostics, got: %v", diagMsg.AsResponse().Error)
}

func TestProjectInfoUntitledFileNotOpened(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	// Test that requesting project info for an untitled file that was never opened
	// does not produce an error - it should return an empty result.
	client := initProjectInfoClient(t, map[string]string{})

	uri := lsproto.DocumentUri("untitled:Untitled-1")

	// custom/projectInfo for a file that was never opened should not error.
	id := client.NextID()
	reqID := lsproto.NewID(lsproto.IntegerOrString{Integer: &id})
	req := lsproto.CustomProjectInfoInfo.NewRequestMessage(reqID, &lsproto.ProjectInfoParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
	})
	resp, ok := client.SendRequestWorker(t, req, reqID)
	assert.Assert(t, ok, "expected a response")
	assert.Assert(t, resp.Error == nil, "expected no error for untitled file that was never opened, got: %v", resp.Error)
}
