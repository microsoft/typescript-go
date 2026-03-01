package lsp_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func initCompletionClient(t *testing.T, files map[string]string, prefs *lsutil.UserPreferences) (*lsptestutil.LSPClient, func() error) {
	t.Helper()

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

	ptrTrue := new(true)
	initMsg, _, ok := lsptestutil.SendRequest(t, client, lsproto.InitializeInfo, &lsproto.InitializeParams{
		Locale: new("en-US"),
		Capabilities: &lsproto.ClientCapabilities{
			TextDocument: &lsproto.TextDocumentClientCapabilities{
				Completion: &lsproto.CompletionClientCapabilities{
					CompletionItem: &lsproto.ClientCompletionItemOptions{
						SnippetSupport:          ptrTrue,
						CommitCharactersSupport: ptrTrue,
					},
					CompletionList: &lsproto.CompletionListCapabilities{
						ItemDefaults: &[]string{"commitCharacters", "editRange"},
					},
				},
			},
		},
	})
	assert.Assert(t, ok && initMsg.AsResponse().Error == nil, "Initialize failed")
	lsptestutil.SendNotification(t, client, lsproto.InitializedInfo, &lsproto.InitializedParams{})
	<-client.Server.InitComplete()

	lsptestutil.SendNotification(t, client, lsproto.WorkspaceDidChangeConfigurationInfo, &lsproto.DidChangeConfigurationParams{
		Settings: map[string]any{"typescript": prefs},
	})

	return client, closeClient
}

func assertCompletionAfterClose(t *testing.T, resp *lsproto.ResponseMessage) {
	t.Helper()
	if resp.Error != nil {
		t.Fatalf("expected no error, got: [%d] %s", resp.Error.Code, resp.Error.Error())
	}
}

func TestAutoImportCompletionAfterFileClose(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	// Close arrives before completion in the dispatch queue.
	// The completion sync phase flushes the queued close, so the captured
	// snapshot has b.ts removed from overlays.
	t.Run("close before completion", func(t *testing.T) {
		t.Parallel()
		prefs := &lsutil.UserPreferences{
			IncludeCompletionsForModuleExports:    core.TSTrue,
			IncludeCompletionsForImportStatements: core.TSTrue,
		}
		client, closeClient := initCompletionClient(t, map[string]string{
			"/home/projects/tsconfig.json": `{"compilerOptions": {"module": "esnext", "target": "esnext"}}`,
			"/home/projects/a.ts":          "export const someVar = 10;",
			"/home/projects/b.ts":          "s",
		}, prefs)
		defer func() {
			if err := closeClient(); err != nil {
				t.Errorf("goroutine error: %v", err)
			}
		}()

		aURI := lsconv.FileNameToDocumentURI("/home/projects/a.ts")
		bURI := lsconv.FileNameToDocumentURI("/home/projects/b.ts")
		lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
			TextDocument: &lsproto.TextDocumentItem{Uri: aURI, LanguageId: "typescript", Text: "export const someVar = 10;"},
		})
		lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
			TextDocument: &lsproto.TextDocumentItem{Uri: bURI, LanguageId: "typescript", Text: "s"},
		})

		lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidCloseInfo, &lsproto.DidCloseTextDocumentParams{
			TextDocument: lsproto.TextDocumentIdentifier{Uri: bURI},
		})

		msg, _, ok := lsptestutil.SendRequest(t, client, lsproto.TextDocumentCompletionInfo, &lsproto.CompletionParams{
			TextDocument: lsproto.TextDocumentIdentifier{Uri: bURI},
			Position:     lsproto.Position{Line: 0, Character: 1},
			Context:      &lsproto.CompletionContext{},
		})
		assert.Assert(t, ok, "expected a response")
		assertCompletionAfterClose(t, msg.AsResponse())
	})

	// Completion is dispatched first; close arrives while the async phase
	// is (likely) in-flight. The main goroutine sends the completion
	// request (guaranteeing it enters the input channel first), while a
	// background goroutine sends the close after a short delay.
	t.Run("close during async", func(t *testing.T) {
		t.Parallel()
		prefs := &lsutil.UserPreferences{
			IncludeCompletionsForModuleExports:    core.TSTrue,
			IncludeCompletionsForImportStatements: core.TSTrue,
		}
		client, closeClient := initCompletionClient(t, map[string]string{
			"/home/projects/tsconfig.json": `{"compilerOptions": {"module": "esnext", "target": "esnext"}}`,
			"/home/projects/a.ts":          "export const someVar = 10;",
			"/home/projects/b.ts":          "s",
		}, prefs)
		defer func() {
			if err := closeClient(); err != nil {
				t.Errorf("goroutine error: %v", err)
			}
		}()

		aURI := lsconv.FileNameToDocumentURI("/home/projects/a.ts")
		bURI := lsconv.FileNameToDocumentURI("/home/projects/b.ts")
		lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
			TextDocument: &lsproto.TextDocumentItem{Uri: aURI, LanguageId: "typescript", Text: "export const someVar = 10;"},
		})
		lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
			TextDocument: &lsproto.TextDocumentItem{Uri: bURI, LanguageId: "typescript", Text: "s"},
		})

		go func() {
			time.Sleep(5 * time.Millisecond)
			lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidCloseInfo, &lsproto.DidCloseTextDocumentParams{
				TextDocument: lsproto.TextDocumentIdentifier{Uri: bURI},
			})
		}()

		msg, _, ok := lsptestutil.SendRequest(t, client, lsproto.TextDocumentCompletionInfo, &lsproto.CompletionParams{
			TextDocument: lsproto.TextDocumentIdentifier{Uri: bURI},
			Position:     lsproto.Position{Line: 0, Character: 1},
			Context:      &lsproto.CompletionContext{},
		})
		assert.Assert(t, ok, "expected a response")
		assertCompletionAfterClose(t, msg.AsResponse())
	})
}

func TestCompletionForUnopenedFile(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	prefs := &lsutil.UserPreferences{}
	client, closeClient := initCompletionClient(t, map[string]string{
		"/home/projects/tsconfig.json": `{"compilerOptions": {"module": "esnext", "target": "esnext"}}`,
		"/home/projects/c.ts":          "let xyz = 1;\n",
	}, prefs)
	defer func() {
		if err := closeClient(); err != nil {
			t.Errorf("goroutine error: %v", err)
		}
	}()

	cURI := lsconv.FileNameToDocumentURI("/home/projects/c.ts")
	msg, _, ok := lsptestutil.SendRequest(t, client, lsproto.TextDocumentCompletionInfo, &lsproto.CompletionParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: cURI},
		Position:     lsproto.Position{Line: 0, Character: 5},
		Context:      &lsproto.CompletionContext{},
	})
	assert.Assert(t, ok, "expected a response")
	resp := msg.AsResponse()
	if resp.Error != nil {
		t.Fatalf("expected no error, got: [%d] %s", resp.Error.Code, resp.Error.Error())
	}
}

func completionItems(resp lsproto.CompletionResponse) []*lsproto.CompletionItem {
	if resp.List != nil {
		return resp.List.Items
	}
	if resp.Items != nil {
		return *resp.Items
	}
	return nil
}

func hasCompletionItem(items []*lsproto.CompletionItem, label string) bool {
	for _, item := range items {
		if item.Label == label {
			return true
		}
	}
	return false
}

// TestCompletionSnapshotFreezing verifies that the auto-import retry uses the
// snapshot captured in the sync phase, not a newer one that includes a
// concurrent DidChange. Without snapshot freezing the retry would flush the
// pending change, making position/prefix inconsistent with the request.
func TestCompletionSnapshotFreezing(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	prefs := &lsutil.UserPreferences{
		IncludeCompletionsForModuleExports:    core.TSTrue,
		IncludeCompletionsForImportStatements: core.TSTrue,
	}
	client, closeClient := initCompletionClient(t, map[string]string{
		"/home/projects/tsconfig.json": `{"compilerOptions": {"module": "esnext", "target": "esnext"}}`,
		"/home/projects/a.ts":          "export const someVar = 10;",
		"/home/projects/b.ts":          "someV",
	}, prefs)
	defer func() {
		if err := closeClient(); err != nil {
			t.Errorf("goroutine error: %v", err)
		}
	}()

	aURI := lsconv.FileNameToDocumentURI("/home/projects/a.ts")
	bURI := lsconv.FileNameToDocumentURI("/home/projects/b.ts")
	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{Uri: aURI, LanguageId: "typescript", Text: "export const someVar = 10;"},
	})
	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{Uri: bURI, LanguageId: "typescript", Text: "someV"},
	})

	go func() {
		time.Sleep(5 * time.Millisecond)
		lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidChangeInfo, &lsproto.DidChangeTextDocumentParams{
			TextDocument: lsproto.VersionedTextDocumentIdentifier{Uri: bURI, Version: 2},
			ContentChanges: []lsproto.TextDocumentContentChangePartialOrWholeDocument{
				{WholeDocument: &lsproto.TextDocumentContentChangeWholeDocument{Text: "notMatching"}},
			},
		})
	}()

	msg, resp, ok := lsptestutil.SendRequest(t, client, lsproto.TextDocumentCompletionInfo, &lsproto.CompletionParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: bURI},
		Position:     lsproto.Position{Line: 0, Character: 5},
		Context:      &lsproto.CompletionContext{},
	})
	assert.Assert(t, ok, "expected a response")
	r := msg.AsResponse()
	if r.Error != nil {
		t.Fatalf("expected no error, got: [%d] %s", r.Error.Code, r.Error.Error())
	}
	assert.Assert(t, hasCompletionItem(completionItems(resp), "someVar"),
		"expected someVar in completions (snapshot freezing should preserve original content)")
}
