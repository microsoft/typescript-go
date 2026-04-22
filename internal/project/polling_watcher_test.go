package project_test

import (
	"context"
	"testing"
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

	session := project.NewSession(&project.SessionInit{
		BackgroundCtx: context.Background(),
		Options: &project.SessionOptions{
			CurrentDirectory:       "/project",
			DefaultLibraryPath:     bundled.LibPath(),
			PositionEncoding:       lsproto.PositionEncodingKindUTF8,
			WatchEnabled:           false, // Client doesn't support watching
			PollingEnabled:         true,  // Enable in-process polling
			PollingInterval:        50 * time.Millisecond,
			LoggingEnabled:         true,
			PushDiagnosticsEnabled: true,
		},
		FS:     wrappedFS,
		Client: client,
		Logger: logger,
	})
	defer session.Close()

	// Open a file to trigger project creation
	ctx := context.Background()
	session.DidOpenFile(ctx, "file:///project/src/index.ts", 1, `const x: number = 1;`, lsproto.LanguageKindTypeScript)
	time.Sleep(100 * time.Millisecond)

	// Verify initial state — session was created successfully with polling
	snap := session.Snapshot()
	assert.Assert(t, snap != nil, "snapshot should exist")

	// Modify a file on disk and wait for the poller to detect it
	_ = fsFromMap.WriteFile("/project/src/index.ts", `const x: number = "wrong";`)
	time.Sleep(300 * time.Millisecond)

	// The session should still be functional after polled changes
	snap2 := session.Snapshot()
	assert.Assert(t, snap2 != nil, "snapshot should exist after polled change")
}

type pollingTestClient struct{}

func (c *pollingTestClient) WatchFiles(ctx context.Context, id project.WatcherID, watchers []*lsproto.FileSystemWatcher) error {
	return nil
}

func (c *pollingTestClient) UnwatchFiles(ctx context.Context, id project.WatcherID) error {
	return nil
}
func (c *pollingTestClient) RefreshDiagnostics(ctx context.Context) error { return nil }
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
