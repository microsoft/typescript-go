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

func initProjectFilesClient(t *testing.T, files map[string]string) *lsptestutil.LSPClient {
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

func TestProjectFilesConfiguredProject(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	client := initProjectFilesClient(t, map[string]string{
		"/home/projects/tsconfig.json": `{}`,
		"/home/projects/index.ts":      "// TODO: implement\nexport const x = 1;",
		"/home/projects/utils.ts":      "export function add(a: number, b: number) { return a + b; }",
	})

	// Open a file to trigger project loading
	uri := lsproto.DocumentUri("file:///home/projects/index.ts")
	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{Uri: uri, LanguageId: "typescript", Text: "// TODO: implement\nexport const x = 1;"},
	})

	msg, resp, ok := lsptestutil.SendRequest(t, client, lsproto.CustomProjectFilesInfo, lsproto.NoParams{})
	assert.Assert(t, ok, "expected a response")
	assert.Assert(t, msg.AsResponse().Error == nil)
	assert.Assert(t, len(resp.Projects) > 0, "expected at least one project")

	// Find the configured project
	var found *lsproto.ProjectFilesProject
	for _, p := range resp.Projects {
		if p.ConfigFilePath == "/home/projects/tsconfig.json" {
			found = p
			break
		}
	}
	assert.Assert(t, found != nil, "expected a configured project with tsconfig.json")
	assert.Assert(t, len(found.Files) >= 2, "expected at least 2 files in the project")

	// Verify both source files are present
	fileSet := make(map[lsproto.DocumentUri]bool)
	for _, f := range found.Files {
		fileSet[f] = true
	}
	assert.Assert(t, fileSet["file:///home/projects/index.ts"], "expected index.ts in project files")
	assert.Assert(t, fileSet["file:///home/projects/utils.ts"], "expected utils.ts in project files")
}

func TestProjectFilesInferredProject(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	client := initProjectFilesClient(t, map[string]string{
		"/home/projects/index.ts": "export const x = 1;",
	})

	uri := lsproto.DocumentUri("file:///home/projects/index.ts")
	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{Uri: uri, LanguageId: "typescript", Text: "export const x = 1;"},
	})

	msg, resp, ok := lsptestutil.SendRequest(t, client, lsproto.CustomProjectFilesInfo, lsproto.NoParams{})
	assert.Assert(t, ok, "expected a response")
	assert.Assert(t, msg.AsResponse().Error == nil)
	assert.Assert(t, len(resp.Projects) > 0, "expected at least one project")

	// Inferred project should have empty configFilePath
	var found *lsproto.ProjectFilesProject
	for _, p := range resp.Projects {
		if p.ConfigFilePath == "" {
			found = p
			break
		}
	}
	assert.Assert(t, found != nil, "expected an inferred project with empty configFilePath")
	assert.Assert(t, len(found.Files) >= 1, "expected at least 1 file in the inferred project")
}

func TestProjectFilesNoProjects(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	// Don't open any files — no projects should be loaded
	client := initProjectFilesClient(t, map[string]string{
		"/home/projects/index.ts": "export const x = 1;",
	})

	msg, resp, ok := lsptestutil.SendRequest(t, client, lsproto.CustomProjectFilesInfo, lsproto.NoParams{})
	assert.Assert(t, ok, "expected a response")
	assert.Assert(t, msg.AsResponse().Error == nil)
	assert.Assert(t, len(resp.Projects) == 0, "expected no projects when no files are opened")
}
