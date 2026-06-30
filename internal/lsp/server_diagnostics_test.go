package lsp_test

import (
	"context"
	"io"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func initDiagnosticClient(t *testing.T, workspaceDiagnosticsCap bool, workspaceDiagnosticsOption bool) (*lsptestutil.LSPClient, lsproto.InitializeResponse, *atomic.Int32, *atomic.Int32) {
	t.Helper()

	fs := bundled.WrapFS(vfstest.FromMap(map[string]string{
		"/home/projects/tsconfig.json": `{"compilerOptions": {"target": "nope"}}`,
		"/home/projects/index.ts":      "export const x = 1;",
	}, false))

	var publishNotifications atomic.Int32
	var refreshRequests atomic.Int32
	onServerRequest := func(_ context.Context, req *lsproto.RequestMessage) *lsproto.ResponseMessage {
		switch req.Method {
		case lsproto.MethodWorkspaceDiagnosticRefresh:
			refreshRequests.Add(1)
			return &lsproto.ResponseMessage{
				ID:      req.ID,
				JSONRPC: req.JSONRPC,
				Result:  lsproto.Null{},
			}
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
	client.OnServerNotification = func(_ context.Context, req *lsproto.RequestMessage) {
		if req.Method == lsproto.MethodTextDocumentPublishDiagnostics {
			publishNotifications.Add(1)
		}
	}

	capabilities := &lsproto.ClientCapabilities{}
	if workspaceDiagnosticsCap {
		capabilities.Workspace = &lsproto.WorkspaceClientCapabilities{
			Diagnostics: &lsproto.DiagnosticWorkspaceClientCapabilities{
				RefreshSupport: new(true),
			},
		}
	}
	initializationOptions := &lsproto.InitializationOptions{
		UseWorkspaceDiagnostics: &workspaceDiagnosticsOption,
	}
	initMsg, initResp, ok := lsptestutil.SendRequest(t, client, lsproto.InitializeInfo, &lsproto.InitializeParams{
		Capabilities: capabilities,
		InitializationOptions: &lsproto.InitializationOptionsOrNull{
			InitializationOptions: initializationOptions,
		},
	})
	assert.Assert(t, ok && initMsg.AsResponse().Error == nil, "Initialize failed")
	lsptestutil.SendNotification(t, client, lsproto.InitializedInfo, &lsproto.InitializedParams{})
	<-client.Server.InitComplete()

	return client, initResp, &publishNotifications, &refreshRequests
}

func TestWorkspaceDiagnosticsCapabilityControlsProgramDiagnosticMode(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	t.Run("workspace diagnostics client pulls tsconfig diagnostics", func(t *testing.T) {
		t.Parallel()
		client, initResp, publishNotifications, refreshRequests := initDiagnosticClient(t, true, true)
		diagnosticProvider := initResp.Capabilities.DiagnosticProvider
		assert.Assert(t, diagnosticProvider != nil && diagnosticProvider.Options != nil, "expected diagnostic provider options")
		assert.Equal(t, diagnosticProvider.Options.WorkspaceDiagnostics, true)

		uri := lsproto.DocumentUri("file:///home/projects/index.ts")
		lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
			TextDocument: &lsproto.TextDocumentItem{Uri: uri, LanguageId: "typescript", Text: "export const x = 1;"},
		})
		docMsg, _, ok := lsptestutil.SendRequest(t, client, lsproto.TextDocumentDiagnosticInfo, &lsproto.DocumentDiagnosticParams{
			TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
		})
		assert.Assert(t, ok && docMsg.AsResponse().Error == nil, "textDocument/diagnostic failed")
		client.Server.Session().WaitForBackgroundTasks()
		msg, report, ok := lsptestutil.SendRequest(t, client, lsproto.WorkspaceDiagnosticInfo, &lsproto.WorkspaceDiagnosticParams{
			PreviousResultIds: []lsproto.PreviousResultId{},
		})
		assert.Assert(t, ok && msg.AsResponse().Error == nil, "workspace/diagnostic failed")

		assert.Equal(t, publishNotifications.Load(), int32(0), "expected no publish diagnostics notifications")
		assert.Assert(t, refreshRequests.Load() > 0, "expected workspace diagnostic refresh request")
		assert.Assert(t, len(report.Items) == 1, "expected one workspace diagnostic item, got: %v", report.Items)
		item := report.Items[0].FullDocumentDiagnosticReport
		assert.Assert(t, item != nil, "expected full workspace diagnostic report")
		assert.Equal(t, item.Uri, lsproto.DocumentUri("file:///home/projects/tsconfig.json"))
		assert.Assert(t, len(item.Items) > 0, "expected tsconfig diagnostics")
		assert.Assert(t, strings.Contains(item.Items[0].Message.AsString(), "Argument for '--target' option must be:"), "unexpected diagnostic: %v", item.Items[0])
	})

	t.Run("legacy client does not advertise workspace diagnostics", func(t *testing.T) {
		t.Parallel()
		_, initResp, _, _ := initDiagnosticClient(t, false, false)
		diagnosticProvider := initResp.Capabilities.DiagnosticProvider
		assert.Assert(t, diagnosticProvider != nil && diagnosticProvider.Options != nil, "expected diagnostic provider options")
		assert.Equal(t, diagnosticProvider.Options.WorkspaceDiagnostics, false)
	})

	t.Run("workspace capability alone does not enable workspace diagnostics", func(t *testing.T) {
		t.Parallel()
		_, initResp, _, _ := initDiagnosticClient(t, true, false)
		diagnosticProvider := initResp.Capabilities.DiagnosticProvider
		assert.Assert(t, diagnosticProvider != nil && diagnosticProvider.Options != nil, "expected diagnostic provider options")
		assert.Equal(t, diagnosticProvider.Options.WorkspaceDiagnostics, false)
	})
}
