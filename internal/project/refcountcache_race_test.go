package project

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project/logging"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

// TestParseCacheConcurrentSnapshotRace reproduces a crash observed in telemetry:
//
//	panic("cache entry not found") at internal/project/refcountcache.go
//	  RefCountCache.Ref -> Project.CreateProgram (single-file clone path)
//	  -> ProjectCollectionBuilder.DidRequestFile -> Snapshot.Clone
//	  -> Session.getSnapshot -> Session.GetLanguageService
//
// The crash is overwhelmingly Windows, occasionally macOS, never Linux, which
// points at a scheduling-sensitive interleaving rather than a logic bug on a
// single goroutine.
//
// Root-cause analysis
// -------------------
// session.parseCache is a RefCountCache keyed by ParseCacheKey{ParseOptions,
// ScriptKind, Hash}. The cache itself is internally synchronized (a
// collections.SyncMap plus a per-entry sync.Mutex), so this is NOT a memory
// data race the -race detector can catch. The panic is a *logical* ref-count
// problem: RefCountCache.Ref(K) loads the entry for K and panics if the entry
// has already been deleted (refCount reached 0 and the entry was removed) by
// the time the load happens.
//
// The single Ref site that takes the panicking path is Project.CreateProgram's
// single-file clone branch, which Refs every reused source file AND every entry
// in newProgram.DuplicateSourceFiles(). The Deref/delete side is
// Snapshot.dispose, which Derefs the same keys when a program's last snapshot
// goes away. The foreground program rebuild is serialized by snapshotUpdateMu,
// but snapshot disposal runs unsynchronized with it (e.g. adoptSnapshotChange's
// discard path and other background Derefs run outside snapshotMu).
//
// DuplicateSourceFiles are the highest-contention entries by far: when the same
// package id (name + version) is installed under several node_modules locations,
// package deduplication keeps one primary copy in SourceFiles() and records the
// rest as DuplicateSourceFiles. Because the copies are byte-identical they all
// map to the SAME ParseCacheKey, so a single clone Refs that one key once per
// duplicate (here: many times) and a single dispose Derefs it just as many
// times. That makes the per-key Ref/Deref traffic on one shared entry an order
// of magnitude heavier than ordinary files, which is exactly the kind of hot
// shared entry the production crash needs to delete out from under a concurrent
// Ref.
//
// This test builds such a deduplicated-package workspace and hammers the window
// from many goroutines for many iterations. It is inherently probabilistic: it
// will not fail on every run, but it is designed to fail "at least some of the
// time", especially under `-count` and on Windows/macOS scheduling:
//
//	go test -run TestParseCacheConcurrentSnapshotRace -count=50 ./internal/project
//
// It is CI-safe: when the race does not trigger it simply passes.
func TestParseCacheConcurrentSnapshotRace(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	const root = "/user/username/projects/myproject"

	// node10 resolution ("module": "commonjs") so a bare `import "x"` resolves to
	// node_modules/x/index.d.ts and picks up the package id from package.json.
	files := map[string]any{
		root + "/tsconfig.json": `{"compilerOptions":{"module":"commonjs","strict":true}}`,
	}

	// Many wrapper packages, each nesting its OWN copy of x@1.0.0. The copies are
	// byte-identical and share package id x@1.0.0, so package deduplication turns
	// all but one into DuplicateSourceFiles that all map to the same parse-cache
	// key. Editing the open src files Refs that single hot key once per duplicate
	// on every clone and Derefs it just as many times on every dispose.
	const wrapperCount = 16
	for i := range wrapperCount {
		wrap := fmt.Sprintf("/node_modules/wrap%d", i)
		files[root+wrap+"/package.json"] = fmt.Sprintf(`{"name":"wrap%d","version":"1.0.0"}`, i)
		files[root+wrap+"/index.d.ts"] = "import X from \"x\";\nexport const w: X;\n"
		files[root+wrap+"/node_modules/x/package.json"] = `{"name":"x","version":"1.0.0"}`
		files[root+wrap+"/node_modules/x/index.d.ts"] = "export default class X {\n    private x: number;\n}\n"
	}

	// A hub file that imports every wrapper, pulling all x copies into the program
	// so the duplicates are formed once for the whole project.
	var hub strings.Builder
	for i := range wrapperCount {
		fmt.Fprintf(&hub, "import { w as w%d } from \"wrap%d\";\n", i, i)
	}
	hub.WriteString("export const hub = 0;\n")
	files[root+"/src/hub.ts"] = hub.String()

	// Several "entry" files that import the hub and are opened/edited concurrently.
	// Editing an open entry takes the single-file clone path in CreateProgram,
	// which Refs all reused files plus all DuplicateSourceFiles.
	const entryCount = 8
	entryContents := make([]string, entryCount)
	for e := range entryCount {
		content := "import { hub } from './hub';\n" + fmt.Sprintf("export const entry%d = hub;\n", e)
		entryContents[e] = content
		files[fmt.Sprintf("%s/src/entry%d.ts", root, e)] = content
	}

	fs := bundled.WrapFS(vfstest.FromMap(files, false /*useCaseSensitiveFileNames*/))
	session := NewSession(&SessionInit{
		BackgroundCtx: context.Background(),
		Logger:        logging.NewTestLogger(),
		Options: &SessionOptions{
			CurrentDirectory:   "/",
			DefaultLibraryPath: bundled.LibPath(),
			TypingsLocation:    "/home/src/Library/Caches/typescript",
			PositionEncoding:   lsproto.PositionEncodingKindUTF8,
			WatchEnabled:       false,
			LoggingEnabled:     false,
		},
		FS: fs,
	})
	defer session.Close()

	ctx := context.Background()
	entryURIs := make([]lsproto.DocumentUri, entryCount)
	for e := range entryCount {
		entryURIs[e] = lsproto.DocumentUri(fmt.Sprintf("file://%s/src/entry%d.ts", root, e))
		session.DidOpenFile(ctx, entryURIs[e], 1, entryContents[e], lsproto.LanguageKindTypeScript)
	}

	// Sanity-check that the workspace actually produced DuplicateSourceFiles; if
	// it didn't, the test isn't exercising the implicated code path at all.
	if _, err := session.GetLanguageService(ctx, entryURIs[0]); err != nil {
		t.Fatalf("GetLanguageService: %v", err)
	}
	dupCount := 0
	for _, project := range session.Snapshot().ProjectCollection.Projects() {
		if project.Program != nil {
			dupCount += len(project.Program.DuplicateSourceFiles())
		}
	}
	if dupCount == 0 {
		t.Fatalf("expected DuplicateSourceFiles to be produced by package deduplication, got 0 (workspace setup is not exercising the duplicate path)")
	}
	t.Logf("DuplicateSourceFiles in initial program: %d", dupCount)

	const iterations = 4000

	// t.Fatal is not safe from non-test goroutines; collect panics and surface
	// them from the main goroutine instead.
	var failMu sync.Mutex
	var failure any
	recordPanic := func() {
		if r := recover(); r != nil {
			failMu.Lock()
			if failure == nil {
				failure = r
			}
			failMu.Unlock()
		}
	}

	var versionMu sync.Mutex
	version := int32(1)
	nextVersion := func() int32 {
		versionMu.Lock()
		defer versionMu.Unlock()
		version++
		return version
	}

	var wg sync.WaitGroup

	// A few editor goroutines that each edit their own open file with a small
	// jitter between edits. The jitter lets the background warmAutoImportCache
	// clone (which DidChangeFile cancels) actually make progress and run its
	// program-clone Ref loop concurrently with foreground update clones, instead
	// of being cancelled immediately by the next edit.
	const editors = 3
	for e := range editors {
		wg.Go(func() {
			defer recordPanic()
			uri := entryURIs[e%entryCount]
			base := entryContents[e%entryCount]
			rng := rand.New(rand.NewPCG(uint64(e), 0x9e3779b9))
			for i := range iterations {
				newText := base + fmt.Sprintf("export const edit_%d_%d = %d;\n", e, i, i)
				session.DidChangeFile(ctx, uri, nextVersion(), []lsproto.TextDocumentContentChangePartialOrWholeDocument{
					{WholeDocument: &lsproto.TextDocumentContentChangeWholeDocument{Text: newText}},
				})
				if rng.IntN(4) == 0 {
					time.Sleep(time.Duration(rng.IntN(200)) * time.Microsecond)
				}
			}
		})
	}

	// Reader goroutines that force foreground snapshot update clones by flushing
	// the editors' pending changes. These foreground CreateProgram Ref loops are
	// the ones that panic when a concurrent background dispose deletes a shared
	// parse-cache entry out from under them.
	const readers = 12
	for r := range readers {
		wg.Go(func() {
			defer recordPanic()
			rng := rand.New(rand.NewPCG(uint64(r), 0x1234567))
			for range iterations {
				uri := entryURIs[rng.IntN(entryCount)]
				if _, err := session.GetLanguageService(ctx, uri); err != nil {
					t.Errorf("GetLanguageService: %v", err)
					return
				}
			}
		})
	}

	// Auto-import requesters that clone the current snapshot WITHOUT holding
	// snapshotMu (the unsynchronized clone path) and let it be adopted or
	// discarded on the background queue. The discard path Derefs (and may dispose)
	// a snapshot outside snapshotMu, racing the foreground Ref loops.
	const autoImporters = 6
	for a := range autoImporters {
		wg.Go(func() {
			defer recordPanic()
			rng := rand.New(rand.NewPCG(uint64(a), 0xdeadbeef))
			for range iterations {
				uri := entryURIs[rng.IntN(entryCount)]
				base := session.Snapshot()
				if !base.tryRef() {
					continue
				}
				if _, err := session.GetLanguageServiceWithAutoImports(ctx, base, uri); err != nil {
					base.Deref(session)
					continue
				}
				base.Deref(session)
			}
		})
	}

	wg.Wait()
	session.WaitForBackgroundTasks()

	failMu.Lock()
	defer failMu.Unlock()
	if failure != nil {
		t.Fatalf("snapshot clone panicked under concurrency (the bug): %v", failure)
	}
}
