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
}
