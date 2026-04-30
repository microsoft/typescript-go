package vfswatch_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"github.com/microsoft/typescript-go/internal/vfs/vfswatch"
)

var defaultDirs = map[string]bool{
	"/src": true, // recursive
}

func newTestFS() vfs.FS {
	return vfstest.FromMap(map[string]string{
		"/src/a.ts":      "const a = 1;",
		"/src/b.ts":      "const b = 2;",
		"/src/c.ts":      "const c = 3;",
		"/src/sub/d.ts":  "const d = 4;",
		"/tsconfig.json": `{}`,
	}, true)
}

func newWatcherWithState(fs vfs.FS) *vfswatch.FileWatcher {
	fw := vfswatch.NewFileWatcher(fs, 10*time.Millisecond, true, func(vfswatch.WatchEvent) {})
	fw.UpdateWatchedDirectories(defaultDirs)
	return fw
}

// TestRaceScanVsUpdateWatchedDirectories tests for data races between
// concurrent ScanForChanges reads and UpdateWatchedDirectories writes
// using recursive directory watching.
func TestRaceScanVsUpdateWatchedDirectories(t *testing.T) {
	t.Parallel()
	fs := newTestFS()
	fw := newWatcherWithState(fs)

	var wg sync.WaitGroup

	for range 10 {
		wg.Go(func() {
			for range 200 {
				fw.ScanForChanges()
			}
		})
	}

	for range 5 {
		wg.Go(func() {
			for range 100 {
				fw.UpdateWatchedDirectories(map[string]bool{"/src": true})
			}
		})
	}

	wg.Wait()
}

// TestRaceNonRecursiveDirectoryAccess tests for data races when
// directories are watched non-recursively, exercising the shallow
// scan path in ScanForChanges while being replaced concurrently.
func TestRaceNonRecursiveDirectoryAccess(t *testing.T) {
	t.Parallel()
	fs := newTestFS()
	fw := vfswatch.NewFileWatcher(fs, 10*time.Millisecond, true, func(vfswatch.WatchEvent) {})
	fw.UpdateWatchedDirectories(map[string]bool{"/src": false})

	var wg sync.WaitGroup

	for range 10 {
		wg.Go(func() {
			for range 200 {
				fw.ScanForChanges()
			}
		})
	}

	for range 5 {
		wg.Go(func() {
			for range 100 {
				fw.UpdateWatchedDirectories(map[string]bool{"/src": false})
			}
		})
	}

	wg.Wait()
}

// TestRacePollIntervalAccess tests for data races on the PollInterval
// field when it is read and written from multiple goroutines.
func TestRacePollIntervalAccess(t *testing.T) {
	t.Parallel()
	fs := newTestFS()
	fw := newWatcherWithState(fs)

	var wg sync.WaitGroup

	for range 10 {
		wg.Go(func() {
			for range 500 {
				fw.ScanForChanges()
			}
		})
	}

	for i := range 5 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := range 200 {
				fw.SetPollInterval(time.Duration(i*200+j) * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
}

// TestRaceMixedOperations hammers all FileWatcher operations
// concurrently: ScanForChanges, UpdateWatchedDirectories, FS mutations,
// and PollInterval writes.
func TestRaceMixedOperations(t *testing.T) {
	t.Parallel()
	fs := newTestFS()
	fw := newWatcherWithState(fs)

	var wg sync.WaitGroup

	// ScanForChanges readers
	for range 8 {
		wg.Go(func() {
			for range 100 {
				fw.ScanForChanges()
			}
		})
	}

	// UpdateWatchedDirectories writers
	for i := range 4 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for range 50 {
				fw.UpdateWatchedDirectories(map[string]bool{"/src": true})
			}
		}(i)
	}

	// FS modifiers
	for i := range 4 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := range 50 {
				path := fmt.Sprintf("/src/gen_%d_%d.ts", i, j)
				_ = fs.WriteFile(path, fmt.Sprintf("const x = %d;", j))
				if j%3 == 0 {
					_ = fs.Remove(path)
				}
			}
		}(i)
	}

	// PollInterval writers
	for i := range 2 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := range 100 {
				fw.SetPollInterval(time.Duration(50+j) * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
}

// TestRaceUpdateWithConcurrentFileModifications creates and deletes
// files on the FS while UpdateWatchedDirectories is scanning the same FS,
// testing for races between the FS walker and concurrent mutations.
func TestRaceUpdateWithConcurrentFileModifications(t *testing.T) {
	t.Parallel()
	fs := newTestFS()
	fw := newWatcherWithState(fs)

	var wg sync.WaitGroup

	// Rapid file creation/deletion
	for i := range 6 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := range 100 {
				path := fmt.Sprintf("/src/churn_%d_%d.ts", i, j)
				_ = fs.WriteFile(path, fmt.Sprintf("export const v = %d;", j))
				_ = fs.Remove(path)
			}
		}(i)
	}

	// Concurrent UpdateWatchedDirectories (walks the FS tree)
	for range 4 {
		wg.Go(func() {
			for range 50 {
				fw.UpdateWatchedDirectories(map[string]bool{"/src": true})
			}
		})
	}

	wg.Wait()
}

// FuzzFileWatcherOperations fuzzes random sequences of file operations
// and watcher state management to find panics and edge cases.
// Run with -race to also detect data races.
func FuzzFileWatcherOperations(f *testing.F) {
	f.Add([]byte{0, 1, 2, 3, 0, 1, 2, 3})
	f.Add([]byte{2, 2, 2, 0, 0, 1, 3, 3})
	f.Add([]byte{3, 3, 3, 3, 0, 0, 0, 0})
	f.Add([]byte{4, 4, 4, 0, 2, 1, 3, 2})
	f.Add([]byte{5, 5, 5, 5, 5, 5, 5, 5})
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0})
	f.Add([]byte{1, 1, 1, 1, 1, 1, 1, 1})

	f.Fuzz(func(t *testing.T, ops []byte) {
		if len(ops) == 0 {
			return
		}

		fs := newTestFS()
		fw := newWatcherWithState(fs)

		files := []string{"/src/a.ts", "/src/b.ts", "/src/c.ts", "/src/new.ts", "/src/sub/new.ts"}

		for i, op := range ops {
			path := files[i%len(files)]

			switch op % 6 {
			case 0: // Write/modify a file
				_ = fs.WriteFile(path, fmt.Sprintf("const x = %d;", i))
			case 1: // Remove a file
				_ = fs.Remove(path)
			case 2: // Scan for changes
				fw.ScanForChanges()
			case 3: // Rebuild watch state
				fw.UpdateWatchedDirectories(defaultDirs)
			case 4: // Update directories and scan
				fw.UpdateWatchedDirectories(map[string]bool{"/src": true})
				fw.ScanForChanges()
			case 5: // Modify PollInterval
				fw.SetPollInterval(time.Duration(i*10) * time.Millisecond)
			}
		}
	})
}

// FuzzFileWatcherConcurrent is a fuzz test that runs random operations
// from multiple goroutines to find concurrency bugs.
func FuzzFileWatcherConcurrent(f *testing.F) {
	f.Add([]byte{0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5})
	f.Add([]byte{0, 0, 0, 3, 3, 3, 2, 2, 2, 1, 1, 1})
	f.Add([]byte{2, 3, 2, 3, 2, 3, 0, 0, 0, 0, 0, 0})

	f.Fuzz(func(t *testing.T, ops []byte) {
		if len(ops) < 4 {
			return
		}

		fs := newTestFS()
		fw := newWatcherWithState(fs)

		files := []string{"/src/a.ts", "/src/b.ts", "/src/c.ts", "/src/new.ts"}

		// Split ops into chunks for different goroutines
		chunkSize := len(ops) / 2
		if chunkSize == 0 {
			chunkSize = 1
		}

		var wg sync.WaitGroup

		for start := 0; start < len(ops); start += chunkSize {
			end := min(start+chunkSize, len(ops))
			chunk := ops[start:end]

			wg.Add(1)
			go func(chunk []byte, goroutineID int) {
				defer wg.Done()
				for i, op := range chunk {
					path := files[(goroutineID*len(chunk)+i)%len(files)]
					switch op % 4 {
					case 0:
						_ = fs.WriteFile(path, fmt.Sprintf("const g%d = %d;", goroutineID, i))
					case 1:
						_ = fs.Remove(path)
					case 2:
						fw.ScanForChanges()
					case 3:
						fw.UpdateWatchedDirectories(map[string]bool{"/src": true})
					}
				}
			}(chunk, start/chunkSize)
		}

		wg.Wait()
	})
}
