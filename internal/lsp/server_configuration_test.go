package lsp_test

import (
	"context"
	"io"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func TestDidChangeWorkspaceConfigurationRequestsCurrentSettings(t *testing.T) {
	t.Parallel()

	editorSettings := map[string]any{
		"tabSize":      4,
		"insertSpaces": true,
	}
	configurationRequests := 0
	onServerRequest := func(_ context.Context, req *lsproto.RequestMessage) *lsproto.ResponseMessage {
		switch req.Method {
		case lsproto.MethodWorkspaceConfiguration:
			configurationRequests++
			params, err := lsproto.UnmarshalParams[*lsproto.ConfigurationParams](req)
			assert.NilError(t, err)
			results := make([]any, len(params.Items))
			for i, item := range params.Items {
				if item.Section != nil && *item.Section == "editor" {
					results[i] = editorSettings
				}
			}
			return &lsproto.ResponseMessage{ID: req.ID, JSONRPC: req.JSONRPC, Result: results}
		case lsproto.MethodClientRegisterCapability:
			return &lsproto.ResponseMessage{ID: req.ID, JSONRPC: req.JSONRPC, Result: lsproto.Null{}}
		default:
			return nil
		}
	}

	client, closeClient := lsptestutil.NewLSPClient(t, lsp.ServerOptions{
		Err:                io.Discard,
		Cwd:                "/home/projects",
		FS:                 bundled.WrapFS(vfstest.FromMap(map[string]string{}, false)),
		DefaultLibraryPath: bundled.LibPath(),
	}, onServerRequest)
	t.Cleanup(func() { _ = closeClient() })

	initMsg, _, ok := lsptestutil.SendRequest(t, client, lsproto.InitializeInfo, &lsproto.InitializeParams{
		Capabilities: &lsproto.ClientCapabilities{
			Workspace: &lsproto.WorkspaceClientCapabilities{Configuration: new(true)},
		},
	})
	assert.Assert(t, ok && initMsg.AsResponse().Error == nil, "Initialize failed")
	lsptestutil.SendNotification(t, client, lsproto.InitializedInfo, &lsproto.InitializedParams{})
	<-client.Server.InitComplete()

	initial := client.Server.Session().Config().FormatCodeSettings
	assert.Equal(t, initial.TabSize, 4)
	assert.Equal(t, initial.IndentSize, 4)
	assert.Equal(t, initial.ConvertTabsToSpaces, core.TSTrue)

	editorSettings = map[string]any{
		"tabSize":      2,
		"insertSpaces": false,
	}
	lsptestutil.SendNotification(t, client, lsproto.WorkspaceDidChangeConfigurationInfo, &lsproto.DidChangeConfigurationParams{
		Settings: map[string]any{"js/ts": map[string]any{}},
	})
	_, _, _ = lsptestutil.SendRequest(t, client, lsproto.WorkspaceSymbolInfo, &lsproto.WorkspaceSymbolParams{})

	updated := client.Server.Session().Config().FormatCodeSettings
	assert.Equal(t, configurationRequests, 2)
	assert.Equal(t, updated.TabSize, 2)
	assert.Equal(t, updated.IndentSize, 2)
	assert.Equal(t, updated.ConvertTabsToSpaces, core.TSFalse)
}
