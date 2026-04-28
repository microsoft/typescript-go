package project_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"gotest.tools/v3/assert"
)

func TestUpdateWatchTimeoutAndRollback(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	files := map[string]any{
		"/home/projects/TS/p1/tsconfig.json": `{
			"compilerOptions": { "noLib": true, "strict": true }
		}`,
		"/home/projects/TS/p1/src/index.ts": `export const x = 1;`,
	}

	t.Run("slow client triggers timeout and rollback, watch succeeds after tsconfig change", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			init, utils := projecttestutil.GetSessionInitOptions(files, nil, &projecttestutil.TypingsInstallerOptions{})

			// Track WatchFiles calls and which watcher IDs were successfully registered.
			var watchCalls atomic.Int32
			var firstBatchDone atomic.Bool
			var mu sync.Mutex
			var successfulIDs []project.WatcherID
			utils.Client().WatchFilesFunc = func(ctx context.Context, id project.WatcherID, _ []*lsproto.FileSystemWatcher) error {
				watchCalls.Add(1)
				if !firstBatchDone.Load() {
					<-ctx.Done()
					return ctx.Err()
				}
				mu.Lock()
				successfulIDs = append(successfulIDs, id)
				mu.Unlock()
				return nil
			}

			session := project.NewSession(init)
			defer session.Close()

			uri := lsproto.DocumentUri("file:///home/projects/TS/p1/src/index.ts")

			// Open the file: triggers updateWatches which calls WatchFiles for
			// the project's watches. All calls block and time out because the
			// client is slow.
			session.DidOpenFile(context.Background(), uri, 1, files["/home/projects/TS/p1/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

			// Let the background goroutine block on WatchFiles, then advance
			// fake time past the 1s watchRequestTimeout.
			synctest.Wait()
			time.Sleep(2 * time.Second)
			synctest.Wait()

			firstCallCount := watchCalls.Load()
			assert.Assert(t, firstCallCount >= 1, "expected at least one WatchFiles call during initial open, got %d", firstCallCount)

			// No watcher IDs should have been successfully registered.
			mu.Lock()
			assert.Equal(t, len(successfulIDs), 0, "expected no successful watches after timeout")
			mu.Unlock()

			// Allow subsequent WatchFiles calls to succeed.
			firstBatchDone.Store(true)

			// Simulate an external tsconfig.json modification. This queues a
			// pending file-change event. When flushed, the config entry is
			// re-parsed and its rootFilesWatch gets a new identity, causing
			// updateWatches to call WatchFiles for that watcher. Without the
			// rollback, the registry would still have a stale entry from the
			// first failed attempt, and Acquire would return false for the
			// overlapping glob — silently skipping the client registration.
			session.DidChangeWatchedFiles(context.Background(), []*lsproto.FileEvent{
				{
					Uri:  "file:///home/projects/TS/p1/tsconfig.json",
					Type: lsproto.FileChangeTypeChanged,
				},
			})

			// DidChangeWatchedFiles only queues the change. Reopen the file to
			// flush pending changes and trigger UpdateSnapshot → updateWatches.
			session.DidOpenFile(context.Background(), uri, 2, files["/home/projects/TS/p1/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

			synctest.Wait()
			time.Sleep(2 * time.Second)
			synctest.Wait()

			// WatchFiles should have been called again for the re-registered watcher.
			retryCallCount := watchCalls.Load()
			assert.Assert(t, retryCallCount > firstCallCount,
				"expected WatchFiles to be called again after tsconfig change, got %d total calls (was %d after first attempt)",
				retryCallCount, firstCallCount)

			// Verify that at least one watcher was successfully registered.
			mu.Lock()
			assert.Assert(t, len(successfulIDs) > 0,
				"expected at least one watcher to be registered after successful retry")
			mu.Unlock()
		})
	})
}
