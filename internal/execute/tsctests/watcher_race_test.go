package tsctests

import (
	"fmt"
	"sync"
	"testing"

	"github.com/microsoft/typescript-go/internal/execute"
)

// createTestWatcher sets up a minimal project with a tsconfig and
// returns a Watcher ready for concurrent testing, plus the TestSys
// for file manipulation.
func createTestWatcher(t *testing.T) (*execute.Watcher, *TestSys) {
	t.Helper()
	input := &tscInput{
		files: FileMap{
			"/home/src/workspaces/project/a.ts":          `const a: number = 1;`,
			"/home/src/workspaces/project/b.ts":          `import { a } from "./a"; export const b = a;`,
			"/home/src/workspaces/project/tsconfig.json": `{}`,
		},
		commandLineArgs: []string{"--watch"},
	}
	sys := newTestSys(input, false)
	result := execute.CommandLine(sys, []string{"--watch"}, sys)
	if result.Watcher == nil {
		t.Fatal("expected Watcher to be non-nil in watch mode")
	}
	w, ok := result.Watcher.(*execute.Watcher)
	if !ok {
		t.Fatalf("expected *execute.Watcher, got %T", result.Watcher)
	}
	return w, sys
}

// TestWatcherConcurrentDoCycle calls DoCycle from multiple goroutines
// while modifying source files, exposing data races on Watcher fields
// such as configModified, program, config, and the underlying
// FileWatcher state. Run with -race to detect.
func TestWatcherConcurrentDoCycle(t *testing.T) {
	t.Parallel()
	w, sys := createTestWatcher(t)

	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				_ = sys.fsFromFileMap().WriteFile(
					"/home/src/workspaces/project/a.ts",
					fmt.Sprintf("const a: number = %d;", i*10+j),
				)
				w.DoCycle()
			}
		}(i)
	}

	wg.Wait()
}

// TestWatcherDoCycleWithConcurrentStateReads calls DoCycle while
// other goroutines read watcher state through the exported test
// helper methods (HasWatchedFilesChanged, WatchStateLen, etc.).
func TestWatcherDoCycleWithConcurrentStateReads(t *testing.T) {
	t.Parallel()
	w, sys := createTestWatcher(t)

	var wg sync.WaitGroup

	// DoCycle goroutines
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 15; j++ {
				_ = sys.fsFromFileMap().WriteFile(
					"/home/src/workspaces/project/a.ts",
					fmt.Sprintf("const a: number = %d;", i*15+j),
				)
				w.DoCycle()
			}
		}(i)
	}

	// State reader goroutines
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				w.DoCycle()
				w.DoCycle()
				w.DoCycle()
				w.DoCycle()
			}
		}()
	}

	wg.Wait()
}

// TestWatcherConcurrentFileChangesAndDoCycle creates, modifies, and
// deletes files from multiple goroutines while DoCycle runs, testing
// races between FS mutations and watch state updates.
func TestWatcherConcurrentFileChangesAndDoCycle(t *testing.T) {
	t.Parallel()
	w, sys := createTestWatcher(t)

	var wg sync.WaitGroup

	// File creators
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				path := fmt.Sprintf("/home/src/workspaces/project/gen_%d_%d.ts", i, j)
				_ = sys.fsFromFileMap().WriteFile(path, fmt.Sprintf("export const x%d_%d = %d;", i, j, j))
			}
		}(i)
	}

	// File deleters
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 20; j++ {
			_ = sys.fsFromFileMap().Remove(
				fmt.Sprintf("/home/src/workspaces/project/gen_0_%d.ts", j),
			)
		}
	}()

	// DoCycle callers
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				w.DoCycle()
			}
		}()
	}

	wg.Wait()
}

// TestWatcherRapidConfigChanges modifies tsconfig.json rapidly from
// multiple goroutines while DoCycle runs, testing races on
// config-related fields (configModified, configHasErrors,
// configFilePaths, config, extendedConfigCache).
func TestWatcherRapidConfigChanges(t *testing.T) {
	t.Parallel()
	w, sys := createTestWatcher(t)

	var wg sync.WaitGroup

	configs := []string{
		`{}`,
		`{"compilerOptions": {"strict": true}}`,
		`{"compilerOptions": {"target": "ES2020"}}`,
		`{"compilerOptions": {"noEmit": true}}`,
	}

	// Config modifiers + DoCycle
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				_ = sys.fsFromFileMap().WriteFile(
					"/home/src/workspaces/project/tsconfig.json",
					configs[(i+j)%len(configs)],
				)
				w.DoCycle()
			}
		}(i)
	}

	// Concurrent source file modifications
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 15; j++ {
				_ = sys.fsFromFileMap().WriteFile(
					"/home/src/workspaces/project/a.ts",
					fmt.Sprintf("const a: number = %d;", i*15+j),
				)
				w.DoCycle()
			}
		}(i)
	}

	// State readers
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 30; j++ {
				w.DoCycle()
				w.DoCycle()
			}
		}()
	}

	wg.Wait()
}

// TestWatcherConcurrentDoCycleNoChanges calls DoCycle from many
// goroutines when no files have changed, testing the early-return
// path where WatchState is read and HasChanges is called. This path
// reads WatchState and WildcardDirectories without synchronization.
func TestWatcherConcurrentDoCycleNoChanges(t *testing.T) {
	t.Parallel()
	w, _ := createTestWatcher(t)

	var wg sync.WaitGroup

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				w.DoCycle()
			}
		}()
	}

	wg.Wait()
}

// TestWatcherAlternatingModifyAndDoCycle alternates between modifying
// a file and calling DoCycle from different goroutines, creating a
// realistic scenario where the file watcher detects changes mid-cycle.
func TestWatcherAlternatingModifyAndDoCycle(t *testing.T) {
	t.Parallel()
	w, sys := createTestWatcher(t)

	var wg sync.WaitGroup

	// Writer goroutine: continuously modifies files
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 100; j++ {
			_ = sys.fsFromFileMap().WriteFile(
				"/home/src/workspaces/project/a.ts",
				fmt.Sprintf("const a: number = %d;", j),
			)
		}
	}()

	// Multiple DoCycle goroutines
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 25; j++ {
				w.DoCycle()
			}
		}()
	}

	// State reader goroutines
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				w.DoCycle()
			}
		}()
	}

	wg.Wait()
}
