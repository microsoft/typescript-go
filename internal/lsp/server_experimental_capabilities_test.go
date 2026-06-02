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

func TestExperimentalCapabilitiesOnExperimentalObject(t *testing.T) {
	t.Parallel()

	files := map[string]string{
		"/home/projects/tsconfig.json": `{}`,
		"/home/projects/test.ts":       `const x = 1;`,
	}

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

	_, result, ok := lsptestutil.SendRequest(t, client, lsproto.InitializeInfo, &lsproto.InitializeParams{
		Capabilities: &lsproto.ClientCapabilities{},
	})
	assert.Assert(t, ok, "Initialize failed")

	// Verify experimental capabilities are set on the experimental object
	assert.Assert(t, result.Capabilities.Experimental != nil, "Experimental capabilities should be set")
	assert.Assert(t, result.Capabilities.Experimental.CustomSourceDefinitionProvider != nil && *result.Capabilities.Experimental.CustomSourceDefinitionProvider, "customSourceDefinitionProvider should be true under experimental")
	assert.Assert(t, result.Capabilities.Experimental.CustomMultiDocumentHighlightProvider != nil && *result.Capabilities.Experimental.CustomMultiDocumentHighlightProvider, "customMultiDocumentHighlightProvider should be true under experimental")
}
