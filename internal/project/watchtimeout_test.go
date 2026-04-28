package project_test

import (
	"context"
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

	// Two separate projects: p1 will be opened first (with slow client) and
	// closed, then p2 will be opened (with working client). The rollback from
	// p1's failed WatchFiles should leave the registry clean so p2's
	// registration succeeds.
	files := map[string]any{
		"/home/projects/TS/p1/tsconfig.json": `{
			"compilerOptions": { "noLib": true, "strict": true }
		}`,
		"/home/projects/TS/p1/src/index.ts": `export const x = 1;`,
		"/home/projects/TS/p2/tsconfig.json": `{
			"compilerOptions": { "noLib": true, "strict": true }
		}`,
		"/home/projects/TS/p2/src/index.ts": `export const y = 2;`,
	}

	t.Run("slow client triggers timeout and rollback, retry succeeds on next update", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			init, utils := projecttestutil.GetSessionInitOptions(files, nil, &projecttestutil.TypingsInstallerOptions{})

			// Track WatchFiles calls. All calls in the first batch (before
			// firstBatchDone is set) block to simulate an unresponsive client;
			// subsequent calls succeed immediately.
			var watchCalls atomic.Int32
			var firstBatchDone atomic.Bool
			utils.Client().WatchFilesFunc = func(ctx context.Context, _ project.WatcherID, _ []*lsproto.FileSystemWatcher) error {
				watchCalls.Add(1)
				if !firstBatchDone.Load() {
					<-ctx.Done()
					return ctx.Err()
				}
				return nil
			}

			session := project.NewSession(init)
			defer session.Close()

			// Open p1: triggers updateWatches which calls WatchFiles. All calls
			// block and time out because the client is slow.
			session.DidOpenFile(context.Background(), "file:///home/projects/TS/p1/src/index.ts", 1, files["/home/projects/TS/p1/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

			// Let the background goroutine block on WatchFiles, then advance
			// fake time past the 1s timeout.
			synctest.Wait()
			time.Sleep(2 * time.Second)
			synctest.Wait()

			firstCallCount := watchCalls.Load()
			assert.Assert(t, firstCallCount >= 1, "expected at least one WatchFiles call, got %d", firstCallCount)

			// Let subsequent WatchFiles calls succeed.
			firstBatchDone.Store(true)

			// Close p1 and open p2. Closing p1 removes its project from the
			// snapshot; opening p2 adds a new project. The diff produces "removed"
			// entries for p1 and "added" entries for p2, forcing updateWatches to
			// re-register watchers.
			session.DidCloseFile(context.Background(), "file:///home/projects/TS/p1/src/index.ts")
			session.DidOpenFile(context.Background(), "file:///home/projects/TS/p2/src/index.ts", 1, files["/home/projects/TS/p2/src/index.ts"].(string), lsproto.LanguageKindTypeScript)

			synctest.Wait()
			time.Sleep(2 * time.Second)
			synctest.Wait()

			// p2's WatchFiles calls should have succeeded.
			retryCallCount := watchCalls.Load()
			assert.Assert(t, retryCallCount > firstCallCount,
				"expected WatchFiles to be called again after rollback, got %d total calls (was %d after first attempt)",
				retryCallCount, firstCallCount)
		})
	})
}
