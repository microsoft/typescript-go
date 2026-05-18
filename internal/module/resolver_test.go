package module_test

import (
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

type resolutionHostStub struct {
	fs  vfs.FS
	cwd string
}

func (h *resolutionHostStub) FS() vfs.FS                  { return h.fs }
func (h *resolutionHostStub) GetCurrentDirectory() string { return h.cwd }

// Regression test for https://github.com/microsoft/typescript-go/issues/3526.
//
// Resolving a node_modules import with a trailing slash (e.g. `pkg/`) must
// produce the same result as without one.
func TestResolveModuleNameTrailingSlash(t *testing.T) {
	t.Parallel()

	fs := vfstest.FromMap(map[string]string{
		"/repo/node_modules/pkg/package.json": `{"name":"pkg","main":"main.js","types":"main.d.ts"}`,
		"/repo/node_modules/pkg/main.d.ts":    "export const x: number;",
		"/repo/node_modules/pkg/main.js":      "exports.x = 1;",
		"/repo/src/file.ts":                   "",
	}, true)
	host := &resolutionHostStub{fs: fs, cwd: "/repo"}
	opts := &core.CompilerOptions{
		ModuleResolution: core.ModuleResolutionKindBundler,
		Module:           core.ModuleKindESNext,
		Target:           core.ScriptTargetESNext,
	}
	resolver := module.NewResolver(host, opts, "", "")

	for _, name := range []string{"pkg", "pkg/"} {
		r, _ := resolver.ResolveModuleName(name, "/repo/src/file.ts", core.ModuleKindESNext, nil)
		if !r.IsResolved() {
			t.Errorf("%q failed to resolve", name)
		}
	}
}

// blockingFS wraps a vfs.FS and forces FileExists calls for `targetPath` to
// block on `gate` until released. It also counts how many goroutines are
// waiting at the gate. This is used to deterministically reproduce the
// `package.json` info-cache insert race described in
// https://github.com/microsoft/typescript-go/issues/3526.
type blockingFS struct {
	vfs.FS
	targetPath string
	gate       chan struct{}
	waiting    atomic.Int32
}

func (f *blockingFS) FileExists(path string) bool {
	if path == f.targetPath {
		f.waiting.Add(1)
		<-f.gate
	}
	return f.FS.FileExists(path)
}

// flipFileExistsFS wraps a vfs.FS and returns false for the first
// FileExists call to `targetPath`, then true for the second. Both calls
// block until released via their respective gate channels. It also blocks
// ReadFile for the target path so that the "file doesn't exist" Set
// completes before the "file exists" Set (reproducing the LoadOrStore race).
type flipFileExistsFS struct {
	vfs.FS
	targetPath     string
	callCount      atomic.Int32
	feWaiting      atomic.Int32
	firstGate      chan struct{}
	secondGate     chan struct{}
	readGate       chan struct{}
	readWaiting    atomic.Int32
}

func (f *flipFileExistsFS) FileExists(path string) bool {
	if path == f.targetPath {
		n := f.callCount.Add(1)
		f.feWaiting.Add(1)
		if n == 1 {
			<-f.firstGate
			return false // first caller: simulate "file not yet visible"
		}
		if n == 2 {
			<-f.secondGate
			return f.FS.FileExists(path) // second caller: file is visible
		}
	}
	return f.FS.FileExists(path)
}

func (f *flipFileExistsFS) ReadFile(path string) (string, bool) {
	if path == f.targetPath {
		f.readWaiting.Add(1)
		<-f.readGate
	}
	return f.FS.ReadFile(path)
}

// Regression test for https://github.com/microsoft/typescript-go/issues/3526.
//
// Two goroutines resolve the same package via specifiers that differ only by
// a trailing slash (`pkg` and `pkg/`). A blocking FS holds both at the
// `FileExists` check for `package.json` — *after* each has confirmed a
// `package.json` info-cache miss but *before* either has called `Set`. When
// released, both proceed to `LoadOrStore` and one of them loses. Without the
// fix, the loser receives the winner's `InfoCacheEntry` whose
// `PackageDirectory` doesn't match its own `candidate` (because one spelling
// has a trailing slash and the other doesn't), and
// `loadNodeModuleFromDirectoryWorker`'s `ComparePaths` check skips loading
// the package's `main`/`types`. With no `index.*` present, resolution falls
// through to "unresolved" — the phantom TS2307 the issue describes. This
// test deterministically fails when the fix is reverted.
func TestResolveModuleNameTrailingSlashRace(t *testing.T) {
	t.Parallel()

	const pkgJSONPath = "/repo/node_modules/pkg/package.json"
	files := map[string]string{
		// `types` points at a file that is not discoverable through any
		// fallback path: there is no `index.*` and no `main`. The only way
		// to resolve `pkg` (or `pkg/`) is via the package.json `types` field
		// inside `loadNodeModuleFromDirectoryWorker`, which is exactly the
		// step that the bug skips when `candidate` and
		// `packageInfo.PackageDirectory` mismatch.
		pkgJSONPath: `{"name":"pkg","types":"./typings/index.d.ts"}`,
		"/repo/node_modules/pkg/typings/index.d.ts": "export const x: number;",
		// Distinct containing files so each `ResolveModuleName` call has a
		// unique module-resolution-cache key.
		"/repo/src/a/file.ts": "",
		"/repo/src/b/file.ts": "",
	}
	fs := &blockingFS{
		FS:         vfstest.FromMap(files, true),
		targetPath: pkgJSONPath,
		gate:       make(chan struct{}),
	}
	host := &resolutionHostStub{fs: fs, cwd: "/repo"}
	opts := &core.CompilerOptions{
		ModuleResolution: core.ModuleResolutionKindBundler,
		Module:           core.ModuleKindESNext,
		Target:           core.ScriptTargetESNext,
	}
	resolver := module.NewResolver(host, opts, "", "")

	type result struct {
		name     string
		resolved bool
	}
	results := make(chan result, 2)
	var wg sync.WaitGroup
	for _, name := range []string{"pkg", "pkg/"} {
		containingFile := "/repo/src/a/file.ts"
		if strings.HasSuffix(name, "/") {
			containingFile = "/repo/src/b/file.ts"
		}
		wg.Go(func() {
			r, _ := resolver.ResolveModuleName(name, containingFile, core.ModuleKindESNext, nil)
			results <- result{name, r.IsResolved()}
		})
	}

	// Wait for both goroutines to reach the FileExists gate, guaranteeing
	// both have observed a package.json info-cache miss.
	deadline := time.Now().Add(5 * time.Second)
	for fs.waiting.Load() < 2 {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for both goroutines to reach FileExists gate; got %d", fs.waiting.Load())
		}
		time.Sleep(time.Millisecond)
	}
	close(fs.gate)

	wg.Wait()
	close(results)
	for r := range results {
		if !r.resolved {
			t.Errorf("%q failed to resolve", r.name)
		}
	}
}

// Regression test for https://github.com/microsoft/typescript-go/issues/1290.
//
// Two goroutines resolve `pkg/sub` concurrently. Both miss the package.json
// info-cache for the root package directory. A `flipFileExistsFS` forces the
// first goroutine's `FileExists` to return false (simulating the file not yet
// being visible), so it stores a nil-Contents cache entry. The second
// goroutine's `FileExists` returns true, but its `Set` call (`LoadOrStore`)
// returns the first goroutine's nil-Contents entry. Without the `Exists()`
// guard on the `typesVersions` lookup, `packageInfo.Contents.GetVersionPaths`
// dereferences nil and panics. With the guard the nil-Contents entry is safely
// skipped.
func TestResolveSubpathNilContentsRace(t *testing.T) {
	t.Parallel()

	const rootPkgJSON = "/repo/node_modules/pkg/package.json"
	files := map[string]string{
		rootPkgJSON:                             `{"name":"pkg","version":"1.0.0"}`,
		"/repo/node_modules/pkg/sub/index.d.ts": "export declare const sub: number;",
		"/repo/node_modules/pkg/sub/index.js":   "exports.sub = 1;",
		"/repo/src/a/file.ts":                   "",
		"/repo/src/b/file.ts":                   "",
	}
	fs := &flipFileExistsFS{
		FS:         vfstest.FromMap(files, true),
		targetPath: rootPkgJSON,
		firstGate:  make(chan struct{}),
		secondGate: make(chan struct{}),
		readGate:   make(chan struct{}),
	}
	host := &resolutionHostStub{fs: fs, cwd: "/repo"}
	opts := &core.CompilerOptions{
		ModuleResolution: core.ModuleResolutionKindBundler,
		Module:           core.ModuleKindESNext,
		Target:           core.ScriptTargetESNext,
	}
	resolver := module.NewResolver(host, opts, "", "")

	var panicked atomic.Bool
	var wg sync.WaitGroup
	// Two goroutines both resolve "pkg/sub". Each calls getPackageJsonInfo
	// for the root package directory, reaching FileExists for rootPkgJSON.
	for _, containingFile := range []string{"/repo/src/a/file.ts", "/repo/src/b/file.ts"} {
		wg.Go(func() {
			defer func() {
				if r := recover(); r != nil {
					panicked.Store(true)
				}
			}()
			resolver.ResolveModuleName("pkg/sub", containingFile, core.ModuleKindESNext, nil)
		})
	}

	// Phase 1: Wait for both goroutines to reach FileExists for the root
	// package.json, guaranteeing both have observed a cache miss.
	deadline := time.Now().Add(5 * time.Second)
	for fs.feWaiting.Load() < 2 {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for both goroutines to reach FileExists gate; got %d", fs.feWaiting.Load())
		}
		time.Sleep(time.Millisecond)
	}

	// Phase 2: Release the first FileExists caller (returns false).
	// It enters the "file not found" branch and stores a nil-Contents entry
	// via Set — this is nearly instant (no ReadFile).
	close(fs.firstGate)

	// Phase 3: Release the second FileExists caller (returns true).
	// It proceeds to ReadFile, which we gate separately to ensure the first
	// goroutine's nil-Contents Set has completed.
	close(fs.secondGate)

	// Phase 4: Wait for the second goroutine to reach ReadFile, then release.
	// By this point the first goroutine has stored its nil-Contents entry.
	// The second goroutine's Set (LoadOrStore) will return that stale entry.
	deadline = time.Now().Add(5 * time.Second)
	for fs.readWaiting.Load() < 1 {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for second goroutine to reach ReadFile gate")
		}
		time.Sleep(time.Millisecond)
	}
	close(fs.readGate)

	wg.Wait()
	if panicked.Load() {
		t.Fatal("resolver panicked due to nil Contents dereference in loadModuleFromSpecificNodeModulesDirectory")
	}
}
