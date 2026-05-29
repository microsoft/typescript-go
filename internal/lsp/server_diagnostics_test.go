package lsp

import (
	"context"
	"errors"
	"io"
	"testing"

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

	ctx := context.Background()
	fs := bundled.WrapFS(vfstest.FromMap(map[string]string{
		"/home/projects/tsconfig.json": `{"compilerOptions": {"strict": true}}`,
		"/home/projects/index.ts":      "let value: string = 1;",
	}, false))
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
