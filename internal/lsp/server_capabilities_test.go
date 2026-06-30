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

func TestInitializeAdvertisesTypeScriptSourceActionKinds(t *testing.T) {
	t.Parallel()

	fs := bundled.WrapFS(vfstest.FromMap(map[string]string{}, false))
	onServerRequest := func(_ context.Context, req *lsproto.RequestMessage) *lsproto.ResponseMessage {
		switch req.Method {
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
		FS:                 fs,
		DefaultLibraryPath: bundled.LibPath(),
	}, onServerRequest)
	t.Cleanup(func() { assert.NilError(t, closeClient()) })

	initMsg, result, ok := lsptestutil.SendRequest(t, client, lsproto.InitializeInfo, &lsproto.InitializeParams{
		Capabilities: &lsproto.ClientCapabilities{},
	})
	assert.Assert(t, ok && initMsg.AsResponse().Error == nil, "Initialize failed")

	codeActionProvider := result.Capabilities.CodeActionProvider
	assert.Assert(t, codeActionProvider != nil && codeActionProvider.CodeActionOptions != nil)
	kinds := codeActionProvider.CodeActionOptions.CodeActionKinds
	assert.Assert(t, kinds != nil)

	for _, kind := range []lsproto.CodeActionKind{
		lsproto.CodeActionKindSourceOrganizeImports,
		lsproto.CodeActionKindSourceOrganizeImportsTs,
		lsproto.CodeActionKindSourceRemoveUnusedImports,
		lsproto.CodeActionKindSourceRemoveUnusedImportsTs,
		lsproto.CodeActionKindSourceSortImports,
		lsproto.CodeActionKindSourceSortImportsTs,
		lsproto.CodeActionKindSourceFixAll,
		lsproto.CodeActionKindSourceFixAllTs,
	} {
		assert.Assert(t, containsCodeActionKind(*kinds, kind), "missing code action kind %q", kind)
	}
}

func containsCodeActionKind(kinds []lsproto.CodeActionKind, kind lsproto.CodeActionKind) bool {
	for _, candidate := range kinds {
		if candidate == kind {
			return true
		}
	}
	return false
}
