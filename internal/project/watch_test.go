package project

import (
	"context"
	"testing"
	"testing/synctest"

	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"gotest.tools/v3/assert"
)

// testWatchClient is a minimal Client implementation for watch tests.
// Only WatchFiles and UnwatchFiles are exercised; other methods are no-ops.
type testWatchClient struct {
	watchFunc   func(ctx context.Context, id WatcherID, watchers []*lsproto.FileSystemWatcher) error
	unwatchFunc func(ctx context.Context, id WatcherID) error
}

func (c *testWatchClient) WatchFiles(ctx context.Context, id WatcherID, watchers []*lsproto.FileSystemWatcher) error {
	if c.watchFunc != nil {
		return c.watchFunc(ctx, id, watchers)
	}
	return nil
}

func (c *testWatchClient) UnwatchFiles(ctx context.Context, id WatcherID) error {
	if c.unwatchFunc != nil {
		return c.unwatchFunc(ctx, id)
	}
	return nil
}

func (c *testWatchClient) RefreshDiagnostics(context.Context) error { return nil }
func (c *testWatchClient) PublishDiagnostics(context.Context, *lsproto.PublishDiagnosticsParams) error {
	return nil
}
func (c *testWatchClient) RefreshInlayHints(context.Context) error                     { return nil }
func (c *testWatchClient) RefreshCodeLens(context.Context) error                       { return nil }
func (c *testWatchClient) ProgressStart(*diagnostics.Message, ...any)                  {}
func (c *testWatchClient) ProgressFinish(*diagnostics.Message, ...any)                 {}
func (c *testWatchClient) SendTelemetry(context.Context, lsproto.TelemetryEvent) error { return nil }
func (c *testWatchClient) IsActive() bool                                              { return true }

func TestGetPathComponentsForWatching(t *testing.T) {
	t.Parallel()

	assert.DeepEqual(t, getPathComponentsForWatching("/project", ""), []string{"/", "project"})
	assert.DeepEqual(t, getPathComponentsForWatching("C:\\project", ""), []string{"C:/", "project"})
	assert.DeepEqual(t, getPathComponentsForWatching("//server/share/project/tsconfig.json", ""), []string{"//server/share", "project", "tsconfig.json"})
	assert.DeepEqual(t, getPathComponentsForWatching(`\\server\share\project\tsconfig.json`, ""), []string{"//server/share", "project", "tsconfig.json"})
	assert.DeepEqual(t, getPathComponentsForWatching("C:\\Users", ""), []string{"C:/Users"})
	assert.DeepEqual(t, getPathComponentsForWatching("C:\\Users\\andrew\\project", ""), []string{"C:/Users/andrew", "project"})
	assert.DeepEqual(t, getPathComponentsForWatching("/home", ""), []string{"/home"})
	assert.DeepEqual(t, getPathComponentsForWatching("/home/andrew/project", ""), []string{"/home/andrew", "project"})
}

func TestNilWatchedFilesClone(t *testing.T) {
	t.Parallel()

	var w *WatchedFiles[int]
	result := w.Clone(42)
	assert.Assert(t, result == nil, "clone on a nil `WatchedFiles` should return nil")
}

func TestUpdateWatchTimeoutAndRollback(t *testing.T) {
	t.Parallel()

	kind := lsproto.WatchKindCreate | lsproto.WatchKindChange | lsproto.WatchKindDelete
	makeWatchedFiles := func() *WatchedFiles[int] {
		return NewWatchedFiles("test", kind, false, func(_ int) PatternsAndIgnored {
			return PatternsAndIgnored{
				patternsInsideWorkspace: []string{"/project/**/*"},
			}
		})
	}

	t.Run("slow client triggers timeout and rollback", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			client := &testWatchClient{
				watchFunc: func(ctx context.Context, _ WatcherID, _ []*lsproto.FileSystemWatcher) error {
					// Simulate an unresponsive client that blocks until the context expires.
					<-ctx.Done()
					return ctx.Err()
				},
			}
			session := &Session{
				watches: newWatchRegistry(),
				client:  client,
			}

			errs := updateWatch(context.Background(), session, nil, nil, makeWatchedFiles())
			assert.Assert(t, len(errs) > 0, "expected timeout errors from slow client")

			// Registry should be empty after rollback.
			assert.Equal(t, len(session.watches.entries), 0, "expected empty registry after rollback")
		})
	})

	t.Run("retry succeeds after rollback", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			calls := 0
			client := &testWatchClient{
				watchFunc: func(ctx context.Context, _ WatcherID, _ []*lsproto.FileSystemWatcher) error {
					calls++
					if calls == 1 {
						// First call: unresponsive client.
						<-ctx.Done()
						return ctx.Err()
					}
					// Subsequent calls: respond immediately.
					return nil
				},
			}
			session := &Session{
				watches: newWatchRegistry(),
				client:  client,
			}

			// First attempt times out and rolls back.
			errs := updateWatch(context.Background(), session, nil, nil, makeWatchedFiles())
			assert.Assert(t, len(errs) > 0, "expected timeout on first attempt")
			assert.Equal(t, len(session.watches.entries), 0, "expected empty registry after rollback")

			// Second attempt (simulates next snapshot update) succeeds.
			errs = updateWatch(context.Background(), session, nil, nil, makeWatchedFiles())
			assert.Assert(t, len(errs) == 0, "expected no errors on retry, got %v", errs)
			assert.Assert(t, len(session.watches.entries) > 0, "expected registry entries after successful retry")
		})
	})
}
