package project_test

import (
	"context"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/project/logging"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func TestPollingWatcherIntegration(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	synctest.Test(t, func(t *testing.T) {
		files := map[string]any{
			"/project/tsconfig.json": `{
				"compilerOptions": {
					"noLib": true,
					"strict": true
				},
				"include": ["src"]
			}`,
			"/project/src/index.ts": `const x: number = 1;`,
			"/project/src/utils.ts": `export function add(a: number, b: number): number { return a + b; }`,
		}

		fsFromMap := vfstest.FromMap(files, false)
		wrappedFS := bundled.WrapFS(fsFromMap)
		logger := logging.NewTestLogger()

		client := &pollingTestClient{}

		bgCtx, bgCancel := context.WithCancel(context.Background())
		defer bgCancel()

		session := project.NewSession(&project.SessionInit{
			BackgroundCtx: bgCtx,
			Options: &project.SessionOptions{
				CurrentDirectory:       "/project",
				DefaultLibraryPath:     bundled.LibPath(),
				PositionEncoding:       lsproto.PositionEncodingKindUTF8,
				WatchEnabled:           false,
				PollingEnabled:         true,
				PollingInterval:        50 * time.Millisecond,
				DebounceDelay:          50 * time.Millisecond,
				LoggingEnabled:         true,
				PushDiagnosticsEnabled: true,
			},
			FS:     wrappedFS,
			Client: client,
			Logger: logger,
		})
		defer session.Close()

		ctx := context.Background()
		session.DidOpenFile(ctx, "file:///project/src/index.ts", 1, `const x: number = 1;`, lsproto.LanguageKindTypeScript)

		// Wait for background tasks (including updatePollingDirectories)
		// to complete so the poller watches the correct directories.
		session.WaitForBackgroundTasks()

		snap := session.Snapshot()
		assert.Assert(t, snap != nil, "snapshot should exist")

		// Advance fake time so the file write gets a different modTime
		// than the snapshot taken by updatePollingDirectories.
		time.Sleep(1 * time.Millisecond)
		synctest.Wait()

		// Modify a file on disk
		_ = fsFromMap.WriteFile("/project/src/index.ts", `const x: number = "wrong";`)

		// Advance fake time past poll interval (50ms) + debounce (250ms) +
		// diagnostics debounce (50ms) and let all goroutines settle.
		time.Sleep(2 * time.Second)
		synctest.Wait()

		assert.Assert(t, client.refreshCount() > 0, "RefreshDiagnostics should have been called after polled file change")

		snap2 := session.Snapshot()
		assert.Assert(t, snap2 != nil, "snapshot should exist after polled change")
	})
}

type pollingTestClient struct {
	mu            sync.Mutex
	refreshCalled int
}

func (c *pollingTestClient) refreshCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.refreshCalled
}

func (c *pollingTestClient) WatchFiles(ctx context.Context, id project.WatcherID, watchers []*lsproto.FileSystemWatcher) error {
	return nil
}

func (c *pollingTestClient) UnwatchFiles(ctx context.Context, id project.WatcherID) error {
	return nil
}

func (c *pollingTestClient) RefreshDiagnostics(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.refreshCalled++
	return nil
}

func (c *pollingTestClient) PublishDiagnostics(ctx context.Context, params *lsproto.PublishDiagnosticsParams) error {
	return nil
}

func (c *pollingTestClient) RefreshInlayHints(ctx context.Context) error              { return nil }
func (c *pollingTestClient) RefreshCodeLens(ctx context.Context) error                { return nil }
func (c *pollingTestClient) ProgressStart(message *diagnostics.Message, args ...any)  {}
func (c *pollingTestClient) ProgressFinish(message *diagnostics.Message, args ...any) {}

func (c *pollingTestClient) SendTelemetry(ctx context.Context, telemetry lsproto.TelemetryEvent) error {
	return nil
}

func (c *pollingTestClient) IsActive() bool { return true }
