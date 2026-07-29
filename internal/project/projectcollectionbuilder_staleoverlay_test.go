package project

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project/logging"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

type staleOverlayTestClient struct{}

var _ Client = (*staleOverlayTestClient)(nil)

func (c *staleOverlayTestClient) WatchFiles(ctx context.Context, id WatcherID, watchers []*lsproto.FileSystemWatcher) error {
	return nil
}
func (c *staleOverlayTestClient) UnwatchFiles(ctx context.Context, id WatcherID) error { return nil }
func (c *staleOverlayTestClient) RefreshDiagnostics(ctx context.Context) error         { return nil }
func (c *staleOverlayTestClient) PublishDiagnostics(ctx context.Context, params *lsproto.PublishDiagnosticsParams) error {
	return nil
}
func (c *staleOverlayTestClient) RefreshInlayHints(ctx context.Context) error              { return nil }
func (c *staleOverlayTestClient) RefreshCodeLens(ctx context.Context) error                { return nil }
func (c *staleOverlayTestClient) ProgressStart(message *diagnostics.Message, args ...any)  {}
func (c *staleOverlayTestClient) ProgressFinish(message *diagnostics.Message, args ...any) {}
func (c *staleOverlayTestClient) SendTelemetry(ctx context.Context, telemetry lsproto.TelemetryEvent) error {
	return nil
}
func (c *staleOverlayTestClient) IsActive() bool           { return false }
func (c *staleOverlayTestClient) SetLocale(locale string)  {}
func (c *staleOverlayTestClient) GetLocale() locale.Locale { return locale.Locale{} }

// TestMarkProjectsAffectedByConfigChangesStaleOverlay reproduces a server crash
// where markProjectsAffectedByConfigChanges dereferences b.fs.overlays[path] for
// a path that is no longer an overlay.
//
// The config file registry caches config file lookups (configFileNames) for open
// files and normally removes an entry when the file's Closed event is processed
// during a snapshot clone. However, the session-level overlay map is updated
// eagerly by flushChangesLocked, *before* the snapshot clone runs. If a snapshot
// update is interrupted between those two steps — most commonly by a panic that
// is recovered at the request boundary (see lsp.Server.recover), which leaves
// the session alive — the pending Closed event has been consumed and the overlay
// removed, but no clone ever performs the registry bookkeeping. The stale
// configFileNames entry then survives indefinitely.
//
// The next config file create/change/delete in an ancestor directory of the
// stale entry (or any excessive-watch-event cache invalidation) reports the
// stale path in changeFileResult.affectedFiles, and the unchecked map lookup in
// markProjectsAffectedByConfigChanges panics inside Session.updateSnapshot —
// while snapshotMu is held, so a host that recovers the panic is left with a
// permanently wedged session.
func TestMarkProjectsAffectedByConfigChangesStaleOverlay(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	files := map[string]string{
		"/src/tsconfig.json": `{"compilerOptions": {"target": "es6"}}`,
		"/src/index.ts":      `export const index = 1;`,
		"/src/other.ts":      `export const other = 2;`,
	}
	testFS := vfstest.FromMap(files, false /*useCaseSensitiveFileNames*/)
	session := NewSession(&SessionInit{
		BackgroundCtx: context.Background(),
		FS:            bundled.WrapFS(testFS),
		Client:        &staleOverlayTestClient{},
		Logger:        logging.NewTestLogger(),
		Options: &SessionOptions{
			CurrentDirectory:   "/",
			DefaultLibraryPath: bundled.LibPath(),
			PositionEncoding:   lsproto.PositionEncodingKindUTF8,
			LoggingEnabled:     true,
		},
	})
	ctx := context.Background()

	session.DidOpenFile(ctx, "file:///src/index.ts", 1, files["/src/index.ts"], lsproto.LanguageKindTypeScript)
	session.DidOpenFile(ctx, "file:///src/other.ts", 1, files["/src/other.ts"], lsproto.LanguageKindTypeScript)

	indexPath := session.toPath("/src/index.ts")
	otherPath := session.toPath("/src/other.ts")

	registry := session.Snapshot().ConfigFileRegistry
	_, ok := registry.configFileNames[indexPath]
	assert.Assert(t, ok, "config lookup for index.ts should be cached after open")
	_, ok = registry.configFileNames[otherPath]
	assert.Assert(t, ok, "config lookup for other.ts should be cached after open")

	// Sever a snapshot update between the overlay-map commit and the registry
	// bookkeeping: flushChangesLocked consumes the pending Closed event and
	// removes the overlay from the session's overlayFS, but the resulting
	// summary is never processed by a snapshot clone. This is the state a
	// session is left in when a panic during a snapshot update is recovered at
	// the request boundary.
	session.pendingFileChangesMu.Lock()
	session.pendingFileChanges = append(session.pendingFileChanges, FileChange{
		Kind: FileChangeKindClose,
		URI:  "file:///src/other.ts",
	})
	_, _ = session.flushChangesLocked(ctx)
	session.pendingFileChangesMu.Unlock()

	// A config file change in an ancestor directory of the stale cached lookup
	// reports it in changeFileResult.affectedFiles.
	assert.NilError(t, testFS.WriteFile("/src/tsconfig.json", `{"compilerOptions": {"target": "esnext"}}`))
	session.DidChangeWatchedFiles(ctx, []*lsproto.FileEvent{
		{
			Uri:  lsproto.DocumentUri("file:///src/tsconfig.json"),
			Type: lsproto.FileChangeTypeChanged,
		},
	})

	// Before the fix, flushing the watch change panicked in
	// markProjectsAffectedByConfigChanges: the stale path has no overlay, and
	// b.fs.overlays[path].FileName() dereferences a nil *Overlay.
	ls, err := session.GetLanguageService(ctx, "file:///src/index.ts")
	assert.NilError(t, err)

	// The project was still reloaded with the changed config, and the open
	// file's default project was recomputed.
	assert.Equal(t, ls.GetProgram().Options().Target, core.ScriptTargetESNext)

	snapshot := session.Snapshot()
	_, ok = snapshot.ConfigFileRegistry.configFileNames[indexPath]
	assert.Assert(t, ok, "config lookup for the open file should be recomputed")
	_, ok = snapshot.ConfigFileRegistry.configFileNames[otherPath]
	assert.Assert(t, !ok, "stale config lookup for the closed file should be dropped")
	_, ok = snapshot.fs.overlays[otherPath]
	assert.Assert(t, !ok, "closed file should not have an overlay")
}
