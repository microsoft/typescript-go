package lsp

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/jsonrpc"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/project/logging"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func TestDocumentDiagnosticCancelsStaleSnapshot(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	server, ctx := newDiagnosticTestServer(t, map[string]string{
		"/home/projects/tsconfig.json": `{"compilerOptions": {"strict": true}}`,
		"/home/projects/index.ts":      "let value: string = 1;",
	})

	uri := lsconv.FileNameToDocumentURI("/home/projects/index.ts")
	server.session.DidOpenFile(ctx, uri, 1, "let value: string = 1;", lsproto.LanguageKindTypeScript)

	requestID := jsonrpc.NewIDInt(1)
	req := lsproto.TextDocumentDiagnosticInfo.NewRequestMessage(requestID, &lsproto.DocumentDiagnosticParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
	})
	doAsyncWork, err := handlers()[lsproto.MethodTextDocumentDiagnostic](server, ctx, req)
	assert.NilError(t, err)
	assert.Assert(t, doAsyncWork != nil)

	server.session.DidChangeFile(ctx, uri, 2, []lsproto.TextDocumentContentChangePartialOrWholeDocument{
		{WholeDocument: &lsproto.TextDocumentContentChangeWholeDocument{Text: `let value: string = "ok";`}},
	})
	_, err = server.session.GetLanguageService(ctx, uri)
	assert.NilError(t, err)

	err = doAsyncWork()
	assert.Assert(t, errors.Is(err, lsproto.ErrorCodeContentModified), "expected stale diagnostics to be rejected, got %v", err)
}

func TestDocumentDiagnosticAllowsUnrelatedPendingChange(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	server, ctx := newDiagnosticTestServer(t, map[string]string{
		"/home/projects/one/tsconfig.json": `{"compilerOptions": {"strict": true}}`,
		"/home/projects/one/index.ts":      "let value: string = 1;",
		"/home/projects/two/tsconfig.json": `{"compilerOptions": {"strict": true}}`,
		"/home/projects/two/other.ts":      "let other: number = 1;",
	})

	uri := lsconv.FileNameToDocumentURI("/home/projects/one/index.ts")
	otherURI := lsconv.FileNameToDocumentURI("/home/projects/two/other.ts")
	server.session.DidOpenFile(ctx, uri, 1, "let value: string = 1;", lsproto.LanguageKindTypeScript)
	server.session.DidOpenFile(ctx, otherURI, 1, "let other: number = 1;", lsproto.LanguageKindTypeScript)

	requestID := jsonrpc.NewIDInt(1)
	req := lsproto.TextDocumentDiagnosticInfo.NewRequestMessage(requestID, &lsproto.DocumentDiagnosticParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
	})
	doAsyncWork, err := handlers()[lsproto.MethodTextDocumentDiagnostic](server, ctx, req)
	assert.NilError(t, err)
	assert.Assert(t, doAsyncWork != nil)

	server.session.DidChangeFile(ctx, otherURI, 2, []lsproto.TextDocumentContentChangePartialOrWholeDocument{
		{WholeDocument: &lsproto.TextDocumentContentChangeWholeDocument{Text: "let other: number = 2;"}},
	})

	err = doAsyncWork()
	assert.NilError(t, err)
	select {
	case msg := <-server.outgoingQueue:
		assert.Assert(t, msg.AsResponse().Error == nil)
	case <-time.After(time.Second):
		t.Fatal("expected diagnostic response")
	}
}

func TestDocumentDiagnosticRejectsMissingParams(t *testing.T) {
	t.Parallel()

	server := NewServer(&ServerOptions{
		Err: io.Discard,
		Cwd: "/home/projects",
		FS:  bundled.WrapFS(vfstest.FromMap(map[string]string{}, false)),
	})
	req := &lsproto.RequestMessage{
		ID:     jsonrpc.NewIDInt(1),
		Method: lsproto.MethodTextDocumentDiagnostic,
	}

	doAsyncWork, err := handlers()[lsproto.MethodTextDocumentDiagnostic](server, context.Background(), req)
	assert.Assert(t, doAsyncWork == nil)
	assert.Assert(t, errors.Is(err, lsproto.ErrorCodeInvalidParams), "expected invalid params, got %v", err)
}

func newDiagnosticTestServer(t *testing.T, files map[string]string) (*Server, context.Context) {
	t.Helper()

	ctx := context.Background()
	fs := bundled.WrapFS(vfstest.FromMap(files, false))
	server := NewServer(&ServerOptions{
		Err:                io.Discard,
		Cwd:                "/home/projects",
		FS:                 fs,
		DefaultLibraryPath: bundled.LibPath(),
	})
	server.backgroundCtx = ctx
	server.positionEncoding = lsproto.PositionEncodingKindUTF16
	server.clientCapabilities = (&lsproto.ClientCapabilities{}).Resolve()
	server.session = project.NewSession(&project.SessionInit{
		BackgroundCtx: lsproto.WithClientCapabilities(ctx, &server.clientCapabilities),
		Options: &project.SessionOptions{
			CurrentDirectory:       "/home/projects",
			DefaultLibraryPath:     bundled.LibPath(),
			PositionEncoding:       lsproto.PositionEncodingKindUTF16,
			PushDiagnosticsEnabled: false,
		},
		FS:     fs,
		Logger: logging.NewTestLogger(),
		Client: server,
	})
	t.Cleanup(server.session.Close)
	server.session.InitializeWithUserConfig(lsutil.NewDefaultUserPreferences())

	return server, ctx
}
