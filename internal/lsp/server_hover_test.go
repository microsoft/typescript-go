package lsp_test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func initHoverClient(t *testing.T, files map[string]string) *lsptestutil.LSPClient {
	t.Helper()

	fs := bundled.WrapFS(vfstest.FromMap(files, false))

	onServerRequest := func(_ context.Context, req *lsproto.RequestMessage) *lsproto.ResponseMessage {
		switch req.Method {
		case lsproto.MethodWorkspaceConfiguration:
			return &lsproto.ResponseMessage{
				ID:      req.ID,
				JSONRPC: req.JSONRPC,
				Result:  []any{&lsutil.UserPreferences{}},
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
	t.Cleanup(func() { _ = closeClient() })

	initMsg, _, ok := lsptestutil.SendRequest(t, client, lsproto.InitializeInfo, &lsproto.InitializeParams{
		Capabilities: &lsproto.ClientCapabilities{},
	})
	assert.Assert(t, ok && initMsg.AsResponse().Error == nil, "Initialize failed")
	lsptestutil.SendNotification(t, client, lsproto.InitializedInfo, &lsproto.InitializedParams{})
	<-client.Server.InitComplete()

	return client
}

// TestHoverAfterFileShrink reproduces the panic from issue #3302 where a
// textDocument/hover request panics with "slice bounds out of range" after a
// file is edited to become much shorter.
func TestHoverAfterFileShrink(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	var longContent strings.Builder
	for i := range 30 {
		longContent.WriteString(fmt.Sprintf("export const var%d = %d;\n", i, i))
	}
	longText := longContent.String()
	shortText := "export const x = 1;\n"

	client := initHoverClient(t, map[string]string{
		"/home/projects/tsconfig.json": `{}`,
		"/home/projects/a.ts":          longText,
	})

	aURI := lsconv.FileNameToDocumentURI("/home/projects/a.ts")

	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{Uri: aURI, LanguageId: "typescript", Text: longText},
	})

	msg, _, ok := lsptestutil.SendRequest(t, client, lsproto.TextDocumentHoverInfo, &lsproto.HoverParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: aURI},
		Position:     lsproto.Position{Line: 0, Character: 15},
	})
	assert.Assert(t, ok, "expected hover response")
	assert.Assert(t, msg.AsResponse().Error == nil, "hover should not error: %v", msg.AsResponse().Error)

	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidChangeInfo, &lsproto.DidChangeTextDocumentParams{
		TextDocument: lsproto.VersionedTextDocumentIdentifier{Uri: aURI, Version: 2},
		ContentChanges: []lsproto.TextDocumentContentChangePartialOrWholeDocument{
			{WholeDocument: &lsproto.TextDocumentContentChangeWholeDocument{Text: shortText}},
		},
	})

	msg, _, ok = lsptestutil.SendRequest(t, client, lsproto.TextDocumentHoverInfo, &lsproto.HoverParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: aURI},
		Position:     lsproto.Position{Line: 0, Character: 15},
	})
	assert.Assert(t, ok, "expected hover response after shrink")
	assert.Assert(t, msg.AsResponse().Error == nil, "hover after shrink should not error: %v", msg.AsResponse().Error)
}

// TestHoverAfterFileShrinkConcurrent tests a race where a hover request is
// dispatched concurrently with a didChange that shrinks the file.
func TestHoverAfterFileShrinkConcurrent(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	var longContent strings.Builder
	for i := range 30 {
		longContent.WriteString(fmt.Sprintf("export const var%d = %d;\n", i, i))
	}
	longText := longContent.String()
	shortText := "export const x = 1;\n"

	client := initHoverClient(t, map[string]string{
		"/home/projects/tsconfig.json": `{}`,
		"/home/projects/a.ts":          longText,
	})

	aURI := lsconv.FileNameToDocumentURI("/home/projects/a.ts")

	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{Uri: aURI, LanguageId: "typescript", Text: longText},
	})

	msg, _, ok := lsptestutil.SendRequest(t, client, lsproto.TextDocumentHoverInfo, &lsproto.HoverParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: aURI},
		Position:     lsproto.Position{Line: 0, Character: 15},
	})
	assert.Assert(t, ok && msg.AsResponse().Error == nil, "initial hover failed")

	waitForHover := lsptestutil.SendRequestAsync(t, client, lsproto.TextDocumentHoverInfo, &lsproto.HoverParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: aURI},
		Position:     lsproto.Position{Line: 5, Character: 10},
	})

	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidChangeInfo, &lsproto.DidChangeTextDocumentParams{
		TextDocument: lsproto.VersionedTextDocumentIdentifier{Uri: aURI, Version: 2},
		ContentChanges: []lsproto.TextDocumentContentChangePartialOrWholeDocument{
			{WholeDocument: &lsproto.TextDocumentContentChangeWholeDocument{Text: shortText}},
		},
	})

	msg, _, ok = waitForHover()
	assert.Assert(t, ok, "expected hover response")
	if msg.AsResponse().Error != nil {
		t.Logf("hover returned error (acceptable): %v", msg.AsResponse().Error)
	}
}

// TestHoverAfterFileShrinkWithNonASCII tests with non-ASCII characters,
// which triggers the UTF-16 position conversion code path that panics in #3302.
func TestHoverAfterFileShrinkWithNonASCII(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	var longContent strings.Builder
	for i := range 30 {
		longContent.WriteString(fmt.Sprintf("export const café%d = \"über value %d\";\n", i, i))
	}
	longText := longContent.String()
	shortText := "export const café = 1;\n"

	client := initHoverClient(t, map[string]string{
		"/home/projects/tsconfig.json": `{}`,
		"/home/projects/a.ts":          longText,
	})

	aURI := lsconv.FileNameToDocumentURI("/home/projects/a.ts")

	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{Uri: aURI, LanguageId: "typescript", Text: longText},
	})

	msg, _, ok := lsptestutil.SendRequest(t, client, lsproto.TextDocumentHoverInfo, &lsproto.HoverParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: aURI},
		Position:     lsproto.Position{Line: 0, Character: 15},
	})
	assert.Assert(t, ok && msg.AsResponse().Error == nil, "initial hover failed")

	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidChangeInfo, &lsproto.DidChangeTextDocumentParams{
		TextDocument: lsproto.VersionedTextDocumentIdentifier{Uri: aURI, Version: 2},
		ContentChanges: []lsproto.TextDocumentContentChangePartialOrWholeDocument{
			{WholeDocument: &lsproto.TextDocumentContentChangeWholeDocument{Text: shortText}},
		},
	})

	msg, _, ok = lsptestutil.SendRequest(t, client, lsproto.TextDocumentHoverInfo, &lsproto.HoverParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: aURI},
		Position:     lsproto.Position{Line: 0, Character: 15},
	})
	assert.Assert(t, ok, "expected hover response after shrink")
	assert.Assert(t, msg.AsResponse().Error == nil, "hover after shrink with non-ASCII should not error: %v", msg.AsResponse().Error)

	waitForHover := lsptestutil.SendRequestAsync(t, client, lsproto.TextDocumentHoverInfo, &lsproto.HoverParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: aURI},
		Position:     lsproto.Position{Line: 0, Character: 10},
	})

	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidChangeInfo, &lsproto.DidChangeTextDocumentParams{
		TextDocument: lsproto.VersionedTextDocumentIdentifier{Uri: aURI, Version: 3},
		ContentChanges: []lsproto.TextDocumentContentChangePartialOrWholeDocument{
			{WholeDocument: &lsproto.TextDocumentContentChangeWholeDocument{Text: "const x = 1;\n"}},
		},
	})

	msg, _, ok = waitForHover()
	assert.Assert(t, ok, "expected hover response")
	if msg.AsResponse().Error != nil {
		t.Logf("concurrent hover returned error (acceptable): %v", msg.AsResponse().Error)
	}
}

// TestHoverAfterPartialEdits tests hover after a series of rapid partial edits
// that shorten the file incrementally, which may exercise different code paths
// than whole-document replacement.
func TestHoverAfterPartialEdits(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	var longContent strings.Builder
	for i := range 30 {
		longContent.WriteString(fmt.Sprintf("export const var%d = %d;\n", i, i))
	}
	longText := longContent.String()

	client := initHoverClient(t, map[string]string{
		"/home/projects/tsconfig.json": `{}`,
		"/home/projects/a.ts":          longText,
	})

	aURI := lsconv.FileNameToDocumentURI("/home/projects/a.ts")

	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{Uri: aURI, LanguageId: "typescript", Text: longText},
	})

	// Initial hover.
	msg, _, ok := lsptestutil.SendRequest(t, client, lsproto.TextDocumentHoverInfo, &lsproto.HoverParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: aURI},
		Position:     lsproto.Position{Line: 0, Character: 15},
	})
	assert.Assert(t, ok && msg.AsResponse().Error == nil, "initial hover failed")

	// Perform multiple rapid partial edits, each deleting a large portion of the file.
	// This replaces lines 1-29 with nothing, leaving only line 0.
	version := int32(2)
	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidChangeInfo, &lsproto.DidChangeTextDocumentParams{
		TextDocument: lsproto.VersionedTextDocumentIdentifier{Uri: aURI, Version: version},
		ContentChanges: []lsproto.TextDocumentContentChangePartialOrWholeDocument{
			{Partial: &lsproto.TextDocumentContentChangePartial{
				Range: lsproto.Range{
					Start: lsproto.Position{Line: 1, Character: 0},
					End:   lsproto.Position{Line: 29, Character: 0},
				},
				Text: "",
			}},
		},
	})

	// Hover immediately after the partial edit.
	msg, _, ok = lsptestutil.SendRequest(t, client, lsproto.TextDocumentHoverInfo, &lsproto.HoverParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: aURI},
		Position:     lsproto.Position{Line: 0, Character: 15},
	})
	assert.Assert(t, ok, "expected hover response after partial edit")
	assert.Assert(t, msg.AsResponse().Error == nil, "hover after partial edit should not error: %v", msg.AsResponse().Error)

	// Do several more rapid edits + hovers to increase chance of hitting a race.
	for i := range 5 {
		version++
		newText := fmt.Sprintf("const v%d = %d;\n", i, i)
		lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidChangeInfo, &lsproto.DidChangeTextDocumentParams{
			TextDocument: lsproto.VersionedTextDocumentIdentifier{Uri: aURI, Version: version},
			ContentChanges: []lsproto.TextDocumentContentChangePartialOrWholeDocument{
				{WholeDocument: &lsproto.TextDocumentContentChangeWholeDocument{Text: newText}},
			},
		})

		msg, _, ok = lsptestutil.SendRequest(t, client, lsproto.TextDocumentHoverInfo, &lsproto.HoverParams{
			TextDocument: lsproto.TextDocumentIdentifier{Uri: aURI},
			Position:     lsproto.Position{Line: 0, Character: 7},
		})
		assert.Assert(t, ok, "expected hover response on iteration %d", i)
		assert.Assert(t, msg.AsResponse().Error == nil, "hover on iteration %d should not error: %v", i, msg.AsResponse().Error)
	}
}
