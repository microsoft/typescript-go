package lsp_test

import (
	"context"
	"io"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func initFileOperationsClient(t *testing.T, files map[string]string, prefs *lsutil.UserPreferences) (*lsptestutil.LSPClient, *lsproto.InitializeResult, vfs.FS) {
	t.Helper()

	if prefs == nil {
		prefs = lsutil.NewDefaultUserPreferences()
	}

	fs := bundled.WrapFS(vfstest.FromMap(files, false))
	onServerRequest := func(_ context.Context, req *lsproto.RequestMessage) *lsproto.ResponseMessage {
		switch req.Method {
		case lsproto.MethodWorkspaceConfiguration:
			return &lsproto.ResponseMessage{
				ID:      req.ID,
				JSONRPC: req.JSONRPC,
				Result:  []any{prefs},
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
		FS:                 fs,
		DefaultLibraryPath: bundled.LibPath(),
	}, onServerRequest)
	t.Cleanup(func() { assert.NilError(t, closeClient()) })

	initMsg, initResp, ok := lsptestutil.SendRequest(t, client, lsproto.InitializeInfo, &lsproto.InitializeParams{
		Capabilities: &lsproto.ClientCapabilities{
			Workspace: &lsproto.WorkspaceClientCapabilities{
				FileOperations: &lsproto.FileOperationClientCapabilities{
					DidRename:  ptr(true),
					WillRename: ptr(true),
				},
			},
		},
	})
	assert.Assert(t, ok && initMsg.AsResponse().Error == nil, "Initialize failed")
	lsptestutil.SendNotification(t, client, lsproto.InitializedInfo, &lsproto.InitializedParams{})
	<-client.Server.InitComplete()

	return client, initResp, fs
}

func ptr[T any](v T) *T {
	return &v
}

func TestInitializeAdvertisesRenameFileOperations(t *testing.T) {
	t.Parallel()

	client, initResp, _ := initFileOperationsClient(t, map[string]string{
		"/home/projects/tsconfig.json": `{"compilerOptions":{"noLib":true},"include":["**/*.ts"]}`,
		"/home/projects/index.ts":      `export const value = 1;`,
	}, nil)
	_ = client

	assert.Assert(t, initResp != nil)
	assert.Assert(t, initResp.Capabilities.Workspace != nil)
	assert.Assert(t, initResp.Capabilities.Workspace.FileOperations != nil)
	assert.Assert(t, initResp.Capabilities.Workspace.FileOperations.WillRename != nil)
	assert.Assert(t, initResp.Capabilities.Workspace.FileOperations.DidRename != nil)

	willRename := initResp.Capabilities.Workspace.FileOperations.WillRename
	assert.Equal(t, len(willRename.Filters), 1)
	assert.Equal(t, *willRename.Filters[0].Scheme, "file")
	assert.Equal(t, willRename.Filters[0].Pattern.Glob, "**/*")
	assert.Equal(t, *willRename.Filters[0].Pattern.Matches, lsproto.FileOperationPatternKindFile)
}
