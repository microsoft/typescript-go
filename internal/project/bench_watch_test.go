package project

import (
	"context"
	"flag"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/fswatch"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/lsp/lspwatcher"
	"github.com/microsoft/typescript-go/internal/project/logging"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

// benchDir can be overridden on the command line, e.g.:
//
// go test ./internal/project/ -bench=BenchmarkWatchRegistration -bench-dir=/path/to/project
//
// When not set it defaults to _submodules/TypeScript (skipped if absent).
var benchDir = flag.String("bench-dir", "", "project root for WatchRegistration benchmark (default: _submodules/TypeScript)")

// benchBackend selects the fswatch backend for the WatchFiles sub-benchmark.
// Valid values: "default", "fsevents", "kqueue", "inotify", "fanotify", "windows".
var (
	benchBackend  = flag.String("bench-backend", "default", "fswatch backend for WatchFiles sub-benchmark (default|fsevents|kqueue|inotify|fanotify|windows)")
	benchGranular = flag.Bool("bench-granular", false, "enable granular watcher computation")
)

// BenchmarkWatchRegistration measures the two separate phases of setting up
// the builtin (in-process) file watcher as driven by the project system.
// Both sub-benchmarks use a real project session so the watcher shapes are
// 100% representative of production (programFilesWatch, typingsWatch,
// rootFilesWatch, config-file watchers, failed-lookup directories, etc.).
//
//   - Watchers: times just the background-queue work that follows DidOpenFile —
//     primarily updateWatches / (WatchedFiles[T]).Watchers() computation — by
//     stopping the timer during session setup and starting it right before
//     WaitForBackgroundTasks.
//   - WatchFiles: replays the captured watcher set against the real fswatch
//     backend. On Linux inotify/fanotify this walks the full tree per root;
//     macOS FSEvents and Windows have cheaper registration.
//
// The sub-benchmarks are kept separate because they are in tension when
// experimenting with granular vs. broad watch strategies: broader patterns
// reduce Watchers computation overhead but increase per-registration cost on
// Linux (larger trees to walk), while more granular patterns do the opposite.
func BenchmarkWatchRegistration(b *testing.B) {
	root := resolveBenchDir(b)
	vfsFS := bundled.WrapFS(osvfs.FS())
	bgCtx := benchClientCapabilities()
	entryURI := benchEntryURI(root)

	// Boot one session upfront to capture the watcher set used by WatchFiles.
	seedClient := &benchClientMock{}
	seedSess := newBenchSession(b, bgCtx, vfsFS, seedClient, root)
	seedSess.DidOpenFile(context.Background(), entryURI, 1, "", lsproto.LanguageKindTypeScript)
	seedSess.WaitForBackgroundTasks()

	allWatchers := seedClient.watchers()
	b.Logf("session registered %d FileSystemWatcher entries", len(allWatchers))
	if len(allWatchers) == 0 {
		b.Fatal("session issued no WatchFiles calls; check that the project was loaded")
	}

	// Sub-benchmark A: Watchers() computation via a real session.
	// Each iteration creates a fresh session, opens the entry file, and waits
	// for the background pass (dominated by updateWatches). Project loading
	// is included in the per-op cost and is approximately constant across
	// strategy changes, so deltas between benchmark runs will accurately
	// reflect changes to watcher computation cost.
	b.Run("Watchers", func(b *testing.B) {
		for b.Loop() {
			client := &benchClientMock{}
			sess := newBenchSession(b, bgCtx, vfsFS, client, root)
			sess.DidOpenFile(context.Background(), entryURI, 1, "", lsproto.LanguageKindTypeScript)
			sess.WaitForBackgroundTasks()
		}
	})

	// Sub-benchmark B: OS watcher registration.
	// Replay the captured watcher list against the real fswatch backend.
	// Close is excluded from timing.
	logger := logging.NewLogger(io.Discard)
	backend := resolveBenchBackend(b)
	b.Logf("fswatch backend: %s", *benchBackend)

	b.Run("WatchFiles", func(b *testing.B) {
		for b.Loop() {
			w := lspwatcher.NewWithFSWatcher(vfsFS, backend, func([]*lsproto.FileEvent) {}, logger)
			if watchErr := w.WatchFiles("bench-id", allWatchers); watchErr != nil {
				b.Fatal(watchErr)
			}
			b.StopTimer()
			w.Close()
			b.StartTimer()
		}
	})
}

func benchClientCapabilities() context.Context {
	caps := lsproto.ResolvedClientCapabilities{
		Workspace: lsproto.ResolvedWorkspaceClientCapabilities{
			DidChangeWatchedFiles: lsproto.ResolvedDidChangeWatchedFilesClientCapabilities{
				DynamicRegistration: true,
			},
		},
	}
	return lsproto.WithClientCapabilities(context.Background(), &caps)
}

func benchEntryURI(root string) lsproto.DocumentUri {
	entry := filepath.Join(root, "src", "compiler", "checker.ts")
	if _, err := os.Stat(entry); err == nil {
		return lsproto.DocumentUri("file://" + tspath.NormalizePath(entry))
	}
	return lsproto.DocumentUri("file://" + tspath.NormalizePath(filepath.Join(root, "compiler", "checker.ts")))
}

func newBenchSession(b *testing.B, bgCtx context.Context, vfsFS vfs.FS, client Client, root string) *Session {
	b.Helper()
	return NewSession(&SessionInit{
		BackgroundCtx: bgCtx,
		FS:            vfsFS,
		Client:        client,
		Logger:        logging.NewLogger(io.Discard),
		Options: &SessionOptions{
			CurrentDirectory:   root,
			DefaultLibraryPath: bundled.LibPath(),
			PositionEncoding:   lsproto.PositionEncodingKindUTF8,
			WatchEnabled:       true,
			GranularWatches:    *benchGranular,
		},
	})
}

// benchClientMock is a minimal project.Client that records WatchFiles calls.
type benchClientMock struct {
	noopClient
	mu  sync.Mutex
	all []*lsproto.FileSystemWatcher
}

func (c *benchClientMock) WatchFiles(_ context.Context, _ WatcherID, ws []*lsproto.FileSystemWatcher) error {
	c.mu.Lock()
	c.all = append(c.all, ws...)
	c.mu.Unlock()
	return nil
}

func (c *benchClientMock) watchers() []*lsproto.FileSystemWatcher {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.all
}

func resolveBenchDir(b *testing.B) string {
	b.Helper()
	if *benchDir != "" {
		return filepath.Clean(*benchDir)
	}
	repo.SkipIfNoTypeScriptSubmodule(b)
	return repo.TypeScriptSubmodulePath()
}

func resolveBenchBackend(b *testing.B) fswatch.Watcher {
	b.Helper()
	switch *benchBackend {
	case "", "default":
		return fswatch.Default()
	case "fsevents":
		return fswatch.FSEvents()
	case "kqueue":
		return fswatch.Kqueue()
	case "inotify":
		return fswatch.Inotify()
	case "fanotify":
		return fswatch.Fanotify()
	case "windows":
		return fswatch.Windows()
	default:
		b.Fatalf("unknown -bench-backend %q; valid values: default, fsevents, kqueue, inotify, fanotify, windows", *benchBackend)
		panic("unreachable")
	}
}
