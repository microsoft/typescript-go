package lsp_test

import (
	"context"
	"io"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
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

func TestWillRenameFilesLoadsProjectsForUnopenedFiles(t *testing.T) {
	t.Parallel()

	client, _, _ := initFileOperationsClient(t, map[string]string{
		"/home/projects/tsconfig.json":  `{"compilerOptions":{"module":"esnext","moduleResolution":"bundler","target":"esnext","noLib":true},"include":["**/*.ts"]}`,
		"/home/projects/shared/util.ts": `export const util = 1;`,
		"/home/projects/foo.ts":         "import { util } from './shared/util';\nexport const foo = util;\n",
		"/home/projects/main.ts":        "import { foo } from \"./foo\";\nfoo;\n",
	}, nil)

	resMsg, result, ok := lsptestutil.SendRequest(t, client, lsproto.WorkspaceWillRenameFilesInfo, &lsproto.RenameFilesParams{
		Files: []*lsproto.FileRename{
			{
				OldUri: string(lsconv.FileNameToDocumentURI("/home/projects/foo.ts")),
				NewUri: string(lsconv.FileNameToDocumentURI("/home/projects/nested/foo.ts")),
			},
		},
	})
	assert.Assert(t, ok && resMsg.AsResponse().Error == nil, "workspace/willRenameFiles failed")
	assert.Assert(t, result.WorkspaceEdit != nil)
	assert.Assert(t, result.WorkspaceEdit.Changes != nil)

	changes := *result.WorkspaceEdit.Changes
	assert.Equal(t, len(changes), 2)

	mainEdits := changes[lsconv.FileNameToDocumentURI("/home/projects/main.ts")]
	assert.Assert(t, len(mainEdits) > 0)
	assert.Equal(t, mainEdits[0].NewText, "\"./nested/foo\"")
}
