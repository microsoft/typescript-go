package vfswatch_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"github.com/microsoft/typescript-go/internal/vfs/vfswatch"
)

func makeLargeFS(nFiles int, nDirs int) vfs.FS {
	files := map[string]string{
		"/tsconfig.json": `{}`,
	}
	for i := range nDirs {
		dir := fmt.Sprintf("/src/dir%d", i)
		for j := range nFiles / nDirs {
			files[fmt.Sprintf("%s/file%d.ts", dir, j)] = fmt.Sprintf("export const x%d_%d = %d;", i, j, i*100+j)
		}
	}
	return vfstest.FromMap(files, true)
}

func makePaths(nFiles int, nDirs int) []string {
	paths := make([]string, 0, nFiles+nDirs+1)
	paths = append(paths, "/tsconfig.json")
	for i := range nDirs {
		dir := fmt.Sprintf("/src/dir%d", i)
		paths = append(paths, dir)
		for j := range nFiles / nDirs {
			paths = append(paths, fmt.Sprintf("%s/file%d.ts", dir, j))
		}
	}
	return paths
}

// BenchmarkUpdateWatchState measures the cost of snapshotting the filesystem.
func BenchmarkUpdateWatchState(b *testing.B) {
	for _, size := range []struct {
		files, dirs int
	}{
		{50, 5},
		{500, 20},
		{2000, 50},
	} {
		b.Run(fmt.Sprintf("files=%d_dirs=%d", size.files, size.dirs), func(b *testing.B) {
			fs := makeLargeFS(size.files, size.dirs)
			fw := vfswatch.NewFileWatcher(fs, 10*time.Millisecond, true, func() {})
			paths := makePaths(size.files, size.dirs)
			wildcards := map[string]bool{"/src": true}

			b.ResetTimer()
			b.ReportAllocs()
			for range b.N {
				fw.UpdateWatchState(paths, wildcards)
			}
		})
	}
}

// BenchmarkHasChangesNoChange measures per-poll cost when nothing changed.
func BenchmarkHasChangesNoChange(b *testing.B) {
	for _, size := range []struct {
		files, dirs int
	}{
		{50, 5},
		{500, 20},
		{2000, 50},
	} {
		b.Run(fmt.Sprintf("files=%d_dirs=%d", size.files, size.dirs), func(b *testing.B) {
			fs := makeLargeFS(size.files, size.dirs)
			fw := vfswatch.NewFileWatcher(fs, 10*time.Millisecond, true, func() {})
			paths := makePaths(size.files, size.dirs)
			wildcards := map[string]bool{"/src": true}
			fw.UpdateWatchState(paths, wildcards)

			b.ResetTimer()
			b.ReportAllocs()
			for range b.N {
				fw.HasChangesFromWatchState()
			}
		})
	}
}

// BenchmarkHasChangesWithChange measures detection cost when one file changed.
func BenchmarkHasChangesWithChange(b *testing.B) {
	for _, size := range []struct {
		files, dirs int
	}{
		{50, 5},
		{500, 20},
		{2000, 50},
	} {
		b.Run(fmt.Sprintf("files=%d_dirs=%d", size.files, size.dirs), func(b *testing.B) {
			fs := makeLargeFS(size.files, size.dirs)
			fw := vfswatch.NewFileWatcher(fs, 10*time.Millisecond, true, func() {})
			paths := makePaths(size.files, size.dirs)
			wildcards := map[string]bool{"/src": true}
			fw.UpdateWatchState(paths, wildcards)
			// Modify one file so hasChanges returns true
			_ = fs.WriteFile("/src/dir0/file0.ts", "export const changed = true;")

			b.ResetTimer()
			b.ReportAllocs()
			for range b.N {
				fw.HasChangesFromWatchState()
			}
		})
	}
}

// BenchmarkHashEntries measures hashing cost for directory listings.
func BenchmarkHashEntries(b *testing.B) {
	for _, nEntries := range []int{10, 100, 500} {
		b.Run(fmt.Sprintf("entries=%d", nEntries), func(b *testing.B) {
			files := map[string]string{}
			for i := range nEntries {
				files[fmt.Sprintf("/dir/file%d.ts", i)] = "x"
			}
			fs := vfstest.FromMap(files, true)
			fw := vfswatch.NewFileWatcher(fs, 10*time.Millisecond, true, func() {})
			paths := []string{"/dir"}
			fw.UpdateWatchState(paths, nil)

			b.ResetTimer()
			b.ReportAllocs()
			for range b.N {
				fw.HasChangesFromWatchState()
			}
		})
	}
}
