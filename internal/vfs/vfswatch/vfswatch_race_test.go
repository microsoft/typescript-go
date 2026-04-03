package vfswatch_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/trackingvfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"github.com/microsoft/typescript-go/internal/vfs/vfswatch"
)

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
	fw := vfswatch.NewFileWatcher(fs, 10*time.Millisecond, true, func() {})
	tfs := &trackingvfs.FS{Inner: fs}
	tfs.SeenFiles.Add("/src/a.ts")
	tfs.SeenFiles.Add("/src/b.ts")
	tfs.SeenFiles.Add("/src/c.ts")
	tfs.SeenFiles.Add("/src/sub/d.ts")
	tfs.SeenFiles.Add("/tsconfig.json")
	fw.UpdateWatchedFiles(tfs)
	return fw
}

// TestRaceHasChangesVsUpdateWatchedFiles tests for data races between
// concurrent HasChanges reads and UpdateWatchedFiles writes on the
// WatchState map. The WatchState field is a plain map with no
// synchronization; concurrent access should be detected by -race.
func TestRaceHasChangesVsUpdateWatchedFiles(t *testing.T) {
	t.Parallel()
	fs := newTestFS()
	fw := newWatcherWithState(fs)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				ws := fw.WatchState
				if ws != nil {
					fw.HasChanges(ws)
				}
			}
		}()
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				tfs := &trackingvfs.FS{Inner: fs}
				tfs.SeenFiles.Add("/src/a.ts")
				tfs.SeenFiles.Add("/src/b.ts")
				fw.UpdateWatchedFiles(tfs)
			}
		}()
	}

	wg.Wait()
}

// TestRaceWildcardDirectoriesAccess tests for data races when
// WildcardDirectories is read internally by HasChanges while being
// replaced concurrently. WildcardDirectories is a plain map assigned
// directly on the struct with no synchronization.
func TestRaceWildcardDirectoriesAccess(t *testing.T) {
	t.Parallel()
	fs := newTestFS()
	fw := newWatcherWithState(fs)
	fw.WildcardDirectories = map[string]bool{"/src": true}

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				ws := fw.WatchState
				if ws != nil {
					fw.HasChanges(ws)
				}
			}
		}()
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				fw.WildcardDirectories = map[string]bool{"/src": true}
			}
		}()
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

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 500; j++ {
				_ = fw.PollInterval
			}
		}()
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				fw.PollInterval = time.Duration(i*200+j) * time.Millisecond
			}
		}(i)
	}

	wg.Wait()
}

// TestRaceMixedOperations hammers all FileWatcher operations
// concurrently: HasChanges, UpdateWatchedFiles, FS mutations,
// WildcardDirectories writes, and PollInterval writes.
func TestRaceMixedOperations(t *testing.T) {
	t.Parallel()
	fs := newTestFS()
	fw := newWatcherWithState(fs)
	fw.WildcardDirectories = map[string]bool{"/src": true}

	var wg sync.WaitGroup

	// HasChanges readers
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				ws := fw.WatchState
				if ws != nil {
					fw.HasChanges(ws)
				}
			}
		}()
	}

	// UpdateWatchedFiles writers
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				tfs := &trackingvfs.FS{Inner: fs}
				tfs.SeenFiles.Add("/src/a.ts")
				tfs.SeenFiles.Add(fmt.Sprintf("/src/new_%d_%d.ts", i, j))
				fw.UpdateWatchedFiles(tfs)
			}
		}(i)
	}

	// FS modifiers
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				path := fmt.Sprintf("/src/gen_%d_%d.ts", i, j)
				_ = fs.WriteFile(path, fmt.Sprintf("const x = %d;", j))
				if j%3 == 0 {
					_ = fs.Remove(path)
				}
			}
		}(i)
	}

	// WildcardDirectories writers
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				fw.WildcardDirectories = map[string]bool{"/src": true}
			}
		}()
	}

	// PollInterval writers
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				fw.PollInterval = time.Duration(50+j) * time.Millisecond
			}
		}(i)
	}

	wg.Wait()
}

// TestRaceUpdateWithConcurrentFileModifications creates and deletes
// files on the FS while UpdateWatchedFiles is scanning the same FS,
// testing for races between the FS walker and concurrent mutations.
func TestRaceUpdateWithConcurrentFileModifications(t *testing.T) {
	t.Parallel()
	fs := newTestFS()
	fw := newWatcherWithState(fs)
	fw.WildcardDirectories = map[string]bool{"/src": true}

	var wg sync.WaitGroup

	// Rapid file creation/deletion
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				path := fmt.Sprintf("/src/churn_%d_%d.ts", i, j)
				_ = fs.WriteFile(path, fmt.Sprintf("export const v = %d;", j))
				_ = fs.Remove(path)
			}
		}(i)
	}

	// Concurrent UpdateWatchedFiles (walks the FS tree via WildcardDirectories)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				tfs := &trackingvfs.FS{Inner: fs}
				tfs.SeenFiles.Add("/src/a.ts")
				tfs.SeenFiles.Add("/tsconfig.json")
				fw.UpdateWatchedFiles(tfs)
			}
		}()
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
			case 2: // Check for changes against current state
				ws := fw.WatchState
				if ws != nil {
					fw.HasChanges(ws)
				}
			case 3: // Rebuild watch state
				tfs := &trackingvfs.FS{Inner: fs}
				for _, f := range files {
					tfs.SeenFiles.Add(f)
				}
				fw.UpdateWatchedFiles(tfs)
			case 4: // Set wildcard directories and check for changes
				fw.WildcardDirectories = map[string]bool{"/src": true}
				ws := fw.WatchState
				if ws != nil {
					fw.HasChanges(ws)
				}
			case 5: // Modify PollInterval
				fw.PollInterval = time.Duration(i*10) * time.Millisecond
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
		fw.WildcardDirectories = map[string]bool{"/src": true}

		files := []string{"/src/a.ts", "/src/b.ts", "/src/c.ts", "/src/new.ts"}

		// Split ops into chunks for different goroutines
		chunkSize := len(ops) / 2
		if chunkSize == 0 {
			chunkSize = 1
		}

		var wg sync.WaitGroup

		for start := 0; start < len(ops); start += chunkSize {
			end := start + chunkSize
			if end > len(ops) {
				end = len(ops)
			}
			chunk := ops[start:end]

			wg.Add(1)
			go func(chunk []byte, goroutineID int) {
				defer wg.Done()
				for i, op := range chunk {
					path := files[(goroutineID*len(chunk)+i)%len(files)]
					switch op % 5 {
					case 0:
						_ = fs.WriteFile(path, fmt.Sprintf("const g%d = %d;", goroutineID, i))
					case 1:
						_ = fs.Remove(path)
					case 2:
						ws := fw.WatchState
						if ws != nil {
							fw.HasChanges(ws)
						}
					case 3:
						tfs := &trackingvfs.FS{Inner: fs}
						tfs.SeenFiles.Add(path)
						fw.UpdateWatchedFiles(tfs)
					case 4:
						fw.WildcardDirectories = map[string]bool{"/src": true}
					}
				}
			}(chunk, start/chunkSize)
		}

		wg.Wait()
	})
}
