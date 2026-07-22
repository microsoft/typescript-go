package lsp_test

import (
	"context"
	"io"
	"sync"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/contentmappertest"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func TestDiscoverContentMappersBeforeDidOpen(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	const component = `<component name="ProfileCard">
<template><h1>{{ title }}</h1></template>
<script lang="ts">
export const title = "Profile";
</script>`
	files := map[string]string{
		"/home/project/tsconfig.json": `{
			"compilerOptions": { "target": "es2020", "module": "esnext", "moduleResolution": "bundler", "strict": true },
			"contentMappers": [ { "package": "mapper", "extensions": [".vue"] } ]
		}`,
		"/home/project/node_modules/mapper/package.json": contentmappertest.PackageJSON(contentmappertest.ComponentMapper),
		"/home/project/ProfileCard.vue":                  component,
	}

	var mu sync.Mutex
	var registrations []*lsproto.Registration
	var unregistrations []*lsproto.Unregistration
	unregisteredSignal := make(chan struct{}, 1)
	onServerRequest := func(_ context.Context, req *lsproto.RequestMessage) *lsproto.ResponseMessage {
		switch req.Method {
		case lsproto.MethodWorkspaceConfiguration:
			return &lsproto.ResponseMessage{ID: req.ID, JSONRPC: req.JSONRPC, Result: []any{nil, nil, nil, nil}}
		case lsproto.MethodClientRegisterCapability:
			params, err := lsproto.UnmarshalParams[*lsproto.RegistrationParams](req)
			assert.NilError(t, err)
			mu.Lock()
			registrations = append(registrations, params.Registrations...)
			mu.Unlock()
			return &lsproto.ResponseMessage{ID: req.ID, JSONRPC: req.JSONRPC, Result: lsproto.Null{}}
		case lsproto.MethodClientUnregisterCapability:
			params, err := lsproto.UnmarshalParams[*lsproto.UnregistrationParams](req)
			assert.NilError(t, err)
			mu.Lock()
			unregistrations = append(unregistrations, params.Unregisterations...)
			mu.Unlock()
			unregisteredSignal <- struct{}{}
			return &lsproto.ResponseMessage{ID: req.ID, JSONRPC: req.JSONRPC, Result: lsproto.Null{}}
		default:
			return nil
		}
	}

	fs := bundled.WrapFS(vfstest.FromMap(files, false))
	client, closeClient := lsptestutil.NewLSPClient(t, lsp.ServerOptions{
		Err:                io.Discard,
		Cwd:                "/home/project",
		FS:                 fs,
		DefaultLibraryPath: bundled.LibPath(),
		Spawn:              contentmappertest.NewSpawner().Spawn,
	}, onServerRequest)
	t.Cleanup(func() { _ = closeClient() })

	caps := &lsproto.ClientCapabilities{TextDocument: &lsproto.TextDocumentClientCapabilities{
		Synchronization: &lsproto.TextDocumentSyncClientCapabilities{DynamicRegistration: new(true)},
	}}
	initMsg, _, ok := lsptestutil.SendRequest(t, client, lsproto.InitializeInfo, &lsproto.InitializeParams{
		Capabilities: caps,
		InitializationOptions: &lsproto.InitializationOptionsOrNull{InitializationOptions: &lsproto.InitializationOptions{
			DangerouslyLoadExternalPlugins: new(true),
		}},
	})
	assert.Assert(t, ok && initMsg.AsResponse().Error == nil, "initialize failed")
	lsptestutil.SendNotification(t, client, lsproto.InitializedInfo, &lsproto.InitializedParams{})
	<-client.Server.InitComplete()

	uri := lsproto.DocumentUri("file:///home/project/ProfileCard.vue")
	msg, result, ok := lsptestutil.SendRequest(t, client, lsproto.CustomDiscoverContentMappersInfo, &lsproto.DiscoverContentMappersParams{
		TextDocuments: []lsproto.TextDocumentIdentifier{{Uri: uri}},
		Extensions:    []string{".vue", ".svelte"},
	})
	assert.Assert(t, ok && msg.AsResponse().Error == nil)
	assert.DeepEqual(t, result.Extensions, []string{".vue"})

	mu.Lock()
	registered := append([]*lsproto.Registration(nil), registrations...)
	mu.Unlock()
	assert.Assert(t, len(registered) > 0, "expected dynamic registrations")
	var foundDidOpen bool
	for _, registration := range registered {
		if registration.Id == "content-mapper-did-open" {
			foundDidOpen = true
			assert.Assert(t, registration.RegisterOptions != nil && registration.RegisterOptions.TextDocumentDidOpen != nil)
			selector := registration.RegisterOptions.TextDocumentDidOpen.DocumentSelector.DocumentSelector
			assert.Assert(t, selector != nil && len(*selector) == 1)
			assert.Equal(t, *(*selector)[0].Pattern.Pattern.Pattern, "**/*.vue")
		}
	}
	assert.Assert(t, foundDidOpen, "expected didOpen registration for .vue")

	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidOpenInfo, &lsproto.DidOpenTextDocumentParams{
		TextDocument: &lsproto.TextDocumentItem{Uri: uri, LanguageId: "vue", Version: 1, Text: component},
	})
	hoverMsg, hover, ok := lsptestutil.SendRequest(t, client, lsproto.TextDocumentHoverInfo, &lsproto.HoverParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
		Position:     lsproto.Position{Line: 3, Character: 15},
	})
	assert.Assert(t, ok && hoverMsg.AsResponse().Error == nil)
	assert.Assert(t, hover.Hover != nil, "expected hover after first foreign didOpen")

	assert.NilError(t, fs.WriteFile("/home/project/tsconfig.json", `{
		"compilerOptions": { "target": "es2020", "module": "esnext", "moduleResolution": "bundler", "strict": true }
	}`))
	lsptestutil.SendNotification(t, client, lsproto.WorkspaceDidChangeWatchedFilesInfo, &lsproto.DidChangeWatchedFilesParams{
		Changes: []*lsproto.FileEvent{{Uri: "file:///home/project/tsconfig.json", Type: lsproto.FileChangeTypeChanged}},
	})
	hoverMsg, hover, _ = lsptestutil.SendRequest(t, client, lsproto.TextDocumentHoverInfo, &lsproto.HoverParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
		Position:     lsproto.Position{Line: 3, Character: 15},
	})
	assert.Assert(t, hoverMsg != nil && hoverMsg.AsResponse().Error == nil, "request before didClose should return a null result")
	assert.Assert(t, hover.Hover == nil)
	diagnosticMsg, diagnostics, ok := lsptestutil.SendRequest(t, client, lsproto.TextDocumentDiagnosticInfo, &lsproto.DocumentDiagnosticParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
	})
	assert.Assert(t, ok && diagnosticMsg.AsResponse().Error == nil, "diagnostics before didClose should return an empty report")
	assert.Assert(t, diagnostics.FullDocumentDiagnosticReport != nil)
	assert.Equal(t, len(diagnostics.FullDocumentDiagnosticReport.Items), 0)
	completionMsg, completion, ok := lsptestutil.SendRequest(t, client, lsproto.TextDocumentCompletionInfo, &lsproto.CompletionParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
		Position:     lsproto.Position{Line: 3, Character: 15},
	})
	assert.Assert(t, ok && completionMsg.AsResponse().Error == nil)
	assert.Assert(t, completion.Items == nil && completion.List == nil)
	referencesMsg, references, ok := lsptestutil.SendRequest(t, client, lsproto.TextDocumentReferencesInfo, &lsproto.ReferenceParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
		Position:     lsproto.Position{Line: 3, Character: 15},
		Context:      &lsproto.ReferenceContext{IncludeDeclaration: true},
	})
	assert.Assert(t, ok && referencesMsg.AsResponse().Error == nil)
	assert.Assert(t, references.Locations == nil)
	renameMsg, rename, ok := lsptestutil.SendRequest(t, client, lsproto.TextDocumentRenameInfo, &lsproto.RenameParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
		Position:     lsproto.Position{Line: 3, Character: 15},
		NewName:      "renamed",
	})
	assert.Assert(t, ok && renameMsg.AsResponse().Error == nil)
	assert.Assert(t, rename.WorkspaceEdit == nil)
	<-unregisteredSignal

	mu.Lock()
	unregistered := append([]*lsproto.Unregistration(nil), unregistrations...)
	mu.Unlock()
	assert.Assert(t, len(unregistered) > 0, "expected dynamic unregistration")
	var foundDidClose bool
	for _, unregistration := range unregistered {
		if unregistration.Id == "content-mapper-did-close" {
			foundDidClose = true
		}
	}
	assert.Assert(t, foundDidClose, "expected didClose unregistration")

	lsptestutil.SendNotification(t, client, lsproto.TextDocumentDidCloseInfo, &lsproto.DidCloseTextDocumentParams{
		TextDocument: lsproto.TextDocumentIdentifier{Uri: uri},
	})
}
