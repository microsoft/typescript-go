package project

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func TestSnapshotFSBuilder(t *testing.T) {
	t.Parallel()

	toPath := func(fileName string) tspath.Path {
		return tspath.Path(fileName)
	}

	t.Run("builds directory tree on file add", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/foo.ts": "const foo = 1;",
		}, false /* useCaseSensitiveFileNames */)

		builder := newSnapshotFSBuilder(
			testFS,
			make(map[tspath.Path]*Overlay), // prevOverlays
			make(map[tspath.Path]*Overlay), // overlays
			make(map[tspath.Path]*diskFile),
			make(map[tspath.Path]*collections.Set[tspath.Path]),
			lsproto.PositionEncodingKindUTF16,
			toPath,
		)

		// Read the file to add it to the diskFiles
		fh := builder.GetFile("/src/foo.ts")
		assert.Assert(t, fh != nil, "file should exist")
		assert.Equal(t, fh.Content(), "const foo = 1;")

		// Finalize and check directories
		snapshot, changed := builder.Finalize()
		assert.Assert(t, changed, "should have changed")

		// Check that directory structure was built
		// /src should contain /src/foo.ts
		srcDir, ok := snapshot.directories[tspath.Path("/src")]
		assert.Assert(t, ok, "/src directory should exist")
		assert.Assert(t, srcDir.Has(tspath.Path("/src/foo.ts")), "/src should contain /src/foo.ts")

		// / should contain /src
		rootDir, ok := snapshot.directories[tspath.Path("/")]
		assert.Assert(t, ok, "/ directory should exist")
		assert.Assert(t, rootDir.Has(tspath.Path("/src")), "/ should contain /src")
	})

	t.Run("builds nested directory tree", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/nested/deep/file.ts": "export const x = 1;",
		}, false /* useCaseSensitiveFileNames */)

		builder := newSnapshotFSBuilder(
			testFS,
			make(map[tspath.Path]*Overlay), // prevOverlays
			make(map[tspath.Path]*Overlay), // overlays
			make(map[tspath.Path]*diskFile),
			make(map[tspath.Path]*collections.Set[tspath.Path]),
			lsproto.PositionEncodingKindUTF16,
			toPath,
		)

		// Read the file to add it to the diskFiles
		fh := builder.GetFile("/src/nested/deep/file.ts")
		assert.Assert(t, fh != nil, "file should exist")

		snapshot, changed := builder.Finalize()
		assert.Assert(t, changed, "should have changed")

		// Check the complete directory tree
		assert.Assert(t, snapshot.directories[tspath.Path("/src/nested/deep")].Has(tspath.Path("/src/nested/deep/file.ts")))
		assert.Assert(t, snapshot.directories[tspath.Path("/src/nested")].Has(tspath.Path("/src/nested/deep")))
		assert.Assert(t, snapshot.directories[tspath.Path("/src")].Has(tspath.Path("/src/nested")))
		assert.Assert(t, snapshot.directories[tspath.Path("/")].Has(tspath.Path("/src")))
	})

	t.Run("removes directory entries on file delete", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/foo.ts": "const foo = 1;",
		}, false /* useCaseSensitiveFileNames */)

		// Start with existing diskFiles and directories
		existingDiskFiles := map[tspath.Path]*diskFile{
			tspath.Path("/src/foo.ts"): newDiskFile("/src/foo.ts", "const foo = 1;"),
		}
		existingDirs := map[tspath.Path]*collections.Set[tspath.Path]{
			tspath.Path("/"):    collections.NewSetFromItems(tspath.Path("/src")),
			tspath.Path("/src"): collections.NewSetFromItems(tspath.Path("/src/foo.ts")),
		}

		builder := newSnapshotFSBuilder(
			testFS,
			make(map[tspath.Path]*Overlay), // prevOverlays
			make(map[tspath.Path]*Overlay), // overlays
			existingDiskFiles,
			existingDirs,
			lsproto.PositionEncodingKindUTF16,
			toPath,
		)

		// Mark the file for deletion by loading and deleting
		if entry, ok := builder.diskFiles.Load(tspath.Path("/src/foo.ts")); ok {
			entry.Delete()
		}

		snapshot, changed := builder.Finalize()
		assert.Assert(t, changed, "should have changed")

		// File should be deleted
		_, hasFile := snapshot.diskFiles[tspath.Path("/src/foo.ts")]
		assert.Assert(t, !hasFile, "file should be deleted")

		// Directory tree should be cleaned up
		_, hasSrcDir := snapshot.directories[tspath.Path("/src")]
		assert.Assert(t, !hasSrcDir, "/src directory should be removed")

		_, hasRootDir := snapshot.directories[tspath.Path("/")]
		assert.Assert(t, !hasRootDir, "root directory should be removed")
	})

	t.Run("removes only empty directories on file delete", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/foo.ts": "const foo = 1;",
			"/src/bar.ts": "const bar = 2;",
		}, false /* useCaseSensitiveFileNames */)

		// Start with existing diskFiles and directories
		existingDiskFiles := map[tspath.Path]*diskFile{
			tspath.Path("/src/foo.ts"): newDiskFile("/src/foo.ts", "const foo = 1;"),
			tspath.Path("/src/bar.ts"): newDiskFile("/src/bar.ts", "const bar = 2;"),
		}
		existingDirs := map[tspath.Path]*collections.Set[tspath.Path]{
			tspath.Path("/"):    collections.NewSetFromItems(tspath.Path("/src")),
			tspath.Path("/src"): collections.NewSetFromItems(tspath.Path("/src/foo.ts"), tspath.Path("/src/bar.ts")),
		}

		builder := newSnapshotFSBuilder(
			testFS,
			make(map[tspath.Path]*Overlay), // prevOverlays
			make(map[tspath.Path]*Overlay), // overlays
			existingDiskFiles,
			existingDirs,
			lsproto.PositionEncodingKindUTF16,
			toPath,
		)

		// Delete only foo.ts
		if entry, ok := builder.diskFiles.Load(tspath.Path("/src/foo.ts")); ok {
			entry.Delete()
		}

		snapshot, changed := builder.Finalize()
		assert.Assert(t, changed, "should have changed")

		// foo.ts should be deleted
		_, hasFile := snapshot.diskFiles[tspath.Path("/src/foo.ts")]
		assert.Assert(t, !hasFile, "foo.ts should be deleted")

		// bar.ts should still exist
		_, hasBar := snapshot.diskFiles[tspath.Path("/src/bar.ts")]
		assert.Assert(t, hasBar, "bar.ts should still exist")

		// /src directory should still exist with bar.ts
		srcDir, hasSrcDir := snapshot.directories[tspath.Path("/src")]
		assert.Assert(t, hasSrcDir, "/src directory should still exist")
		assert.Assert(t, !srcDir.Has(tspath.Path("/src/foo.ts")), "/src should not contain foo.ts")
		assert.Assert(t, srcDir.Has(tspath.Path("/src/bar.ts")), "/src should contain bar.ts")

		// root should still contain /src
		rootDir, hasRootDir := snapshot.directories[tspath.Path("/")]
		assert.Assert(t, hasRootDir, "root directory should still exist")
		assert.Assert(t, rootDir.Has(tspath.Path("/src")), "root should contain /src")
	})

	t.Run("adds file to existing directory", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/foo.ts": "const foo = 1;",
			"/src/bar.ts": "const bar = 2;",
		}, false /* useCaseSensitiveFileNames */)

		// Start with existing file and directories
		existingDiskFiles := map[tspath.Path]*diskFile{
			tspath.Path("/src/foo.ts"): newDiskFile("/src/foo.ts", "const foo = 1;"),
		}
		existingDirs := map[tspath.Path]*collections.Set[tspath.Path]{
			tspath.Path("/"):    collections.NewSetFromItems(tspath.Path("/src")),
			tspath.Path("/src"): collections.NewSetFromItems(tspath.Path("/src/foo.ts")),
		}

		builder := newSnapshotFSBuilder(
			testFS,
			make(map[tspath.Path]*Overlay), // prevOverlays
			make(map[tspath.Path]*Overlay), // overlays
			existingDiskFiles,
			existingDirs,
			lsproto.PositionEncodingKindUTF16,
			toPath,
		)

		// Read bar.ts to add it
		fh := builder.GetFile("/src/bar.ts")
		assert.Assert(t, fh != nil, "bar.ts should exist")

		snapshot, changed := builder.Finalize()
		assert.Assert(t, changed, "should have changed")

		// /src should contain both files
		srcDir := snapshot.directories[tspath.Path("/src")]
		assert.Assert(t, srcDir.Has(tspath.Path("/src/foo.ts")), "/src should contain foo.ts")
		assert.Assert(t, srcDir.Has(tspath.Path("/src/bar.ts")), "/src should contain bar.ts")
	})

	t.Run("no change when no files added or deleted", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/foo.ts": "const foo = 1;",
		}, false /* useCaseSensitiveFileNames */)

		existingDiskFiles := map[tspath.Path]*diskFile{
			tspath.Path("/src/foo.ts"): newDiskFile("/src/foo.ts", "const foo = 1;"),
		}
		existingDirs := map[tspath.Path]*collections.Set[tspath.Path]{
			tspath.Path("/"):    collections.NewSetFromItems(tspath.Path("/src")),
			tspath.Path("/src"): collections.NewSetFromItems(tspath.Path("/src/foo.ts")),
		}

		builder := newSnapshotFSBuilder(
			testFS,
			make(map[tspath.Path]*Overlay), // prevOverlays
			make(map[tspath.Path]*Overlay), // overlays
			existingDiskFiles,
			existingDirs,
			lsproto.PositionEncodingKindUTF16,
			toPath,
		)

		// Don't add or delete any files
		snapshot, changed := builder.Finalize()
		assert.Assert(t, !changed, "should not have changed")

		// Directories should remain the same
		srcDir := snapshot.directories[tspath.Path("/src")]
		assert.Assert(t, srcDir.Has(tspath.Path("/src/foo.ts")))
	})

	t.Run("overlay files are returned over disk files", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/foo.ts": "const foo = 1;",
		}, false /* useCaseSensitiveFileNames */)

		overlays := map[tspath.Path]*Overlay{
			tspath.Path("/src/foo.ts"): {
				fileBase: fileBase{fileName: "/src/foo.ts", content: "const foo = 999;"},
			},
		}

		builder := newSnapshotFSBuilder(
			testFS,
			make(map[tspath.Path]*Overlay), // prevOverlays
			overlays,
			make(map[tspath.Path]*diskFile),
			make(map[tspath.Path]*collections.Set[tspath.Path]),
			lsproto.PositionEncodingKindUTF16,
			toPath,
		)

		// Should return overlay content
		fh := builder.GetFile("/src/foo.ts")
		assert.Assert(t, fh != nil)
		assert.Equal(t, fh.Content(), "const foo = 999;")
	})

	t.Run("multiple files added and deleted in single cycle", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/a.ts":        "const a = 1;",
			"/src/b.ts":        "const b = 2;",
			"/lib/utils.ts":    "export const util = 1;",
			"/lib/helpers.ts":  "export const helper = 1;",
			"/other/single.ts": "const single = 1;",
		}, false /* useCaseSensitiveFileNames */)

		// Start with some existing files
		existingDiskFiles := map[tspath.Path]*diskFile{
			tspath.Path("/src/a.ts"):        newDiskFile("/src/a.ts", "const a = 1;"),
			tspath.Path("/other/single.ts"): newDiskFile("/other/single.ts", "const single = 1;"),
		}
		existingDirs := map[tspath.Path]*collections.Set[tspath.Path]{
			tspath.Path("/"):      collections.NewSetFromItems(tspath.Path("/src"), tspath.Path("/other")),
			tspath.Path("/src"):   collections.NewSetFromItems(tspath.Path("/src/a.ts")),
			tspath.Path("/other"): collections.NewSetFromItems(tspath.Path("/other/single.ts")),
		}

		builder := newSnapshotFSBuilder(
			testFS,
			make(map[tspath.Path]*Overlay), // prevOverlays
			make(map[tspath.Path]*Overlay), // overlays
			existingDiskFiles,
			existingDirs,
			lsproto.PositionEncodingKindUTF16,
			toPath,
		)

		// Add new files
		fh := builder.GetFile("/src/b.ts")
		assert.Assert(t, fh != nil)
		fh = builder.GetFile("/lib/utils.ts")
		assert.Assert(t, fh != nil)
		fh = builder.GetFile("/lib/helpers.ts")
		assert.Assert(t, fh != nil)

		// Delete existing files
		if entry, ok := builder.diskFiles.Load(tspath.Path("/src/a.ts")); ok {
			entry.Delete()
		}
		if entry, ok := builder.diskFiles.Load(tspath.Path("/other/single.ts")); ok {
			entry.Delete()
		}

		snapshot, changed := builder.Finalize()
		assert.Assert(t, changed, "should have changed")

		// Verify deleted files are gone
		_, hasA := snapshot.diskFiles[tspath.Path("/src/a.ts")]
		assert.Assert(t, !hasA, "/src/a.ts should be deleted")
		_, hasSingle := snapshot.diskFiles[tspath.Path("/other/single.ts")]
		assert.Assert(t, !hasSingle, "/other/single.ts should be deleted")

		// Verify added files exist
		_, hasB := snapshot.diskFiles[tspath.Path("/src/b.ts")]
		assert.Assert(t, hasB, "/src/b.ts should exist")
		_, hasUtils := snapshot.diskFiles[tspath.Path("/lib/utils.ts")]
		assert.Assert(t, hasUtils, "/lib/utils.ts should exist")
		_, hasHelpers := snapshot.diskFiles[tspath.Path("/lib/helpers.ts")]
		assert.Assert(t, hasHelpers, "/lib/helpers.ts should exist")

		// Verify /other directory is cleaned up (was only entry deleted)
		_, hasOther := snapshot.directories[tspath.Path("/other")]
		assert.Assert(t, !hasOther, "/other directory should be removed")

		// Verify /src still exists with b.ts (a.ts deleted, b.ts added)
		srcDir, hasSrc := snapshot.directories[tspath.Path("/src")]
		assert.Assert(t, hasSrc, "/src directory should exist")
		assert.Assert(t, !srcDir.Has(tspath.Path("/src/a.ts")), "/src should not contain a.ts")
		assert.Assert(t, srcDir.Has(tspath.Path("/src/b.ts")), "/src should contain b.ts")

		// Verify /lib was created with both files
		libDir, hasLib := snapshot.directories[tspath.Path("/lib")]
		assert.Assert(t, hasLib, "/lib directory should exist")
		assert.Assert(t, libDir.Has(tspath.Path("/lib/utils.ts")), "/lib should contain utils.ts")
		assert.Assert(t, libDir.Has(tspath.Path("/lib/helpers.ts")), "/lib should contain helpers.ts")

		// Verify root contains /src and /lib but not /other
		rootDir := snapshot.directories[tspath.Path("/")]
		assert.Assert(t, rootDir.Has(tspath.Path("/src")), "root should contain /src")
		assert.Assert(t, rootDir.Has(tspath.Path("/lib")), "root should contain /lib")
		assert.Assert(t, !rootDir.Has(tspath.Path("/other")), "root should not contain /other")
	})

	t.Run("new overlay adds to directory structure", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{}, false /* useCaseSensitiveFileNames */)

		// No previous overlays, new overlay is added
		prevOverlays := make(map[tspath.Path]*Overlay)
		overlays := map[tspath.Path]*Overlay{
			tspath.Path("/src/new.ts"): {
				fileBase: fileBase{fileName: "/src/new.ts", content: "const new = 1;"},
			},
		}

		builder := newSnapshotFSBuilder(
			testFS,
			prevOverlays,
			overlays,
			make(map[tspath.Path]*diskFile),
			make(map[tspath.Path]*collections.Set[tspath.Path]),
			lsproto.PositionEncodingKindUTF16,
			toPath,
		)

		snapshot, _ := builder.Finalize()

		// Check that overlay file is in directory structure
		srcDir, ok := snapshot.directories[tspath.Path("/src")]
		assert.Assert(t, ok, "/src directory should exist")
		assert.Assert(t, srcDir.Has(tspath.Path("/src/new.ts")), "/src should contain new.ts")

		rootDir, ok := snapshot.directories[tspath.Path("/")]
		assert.Assert(t, ok, "/ directory should exist")
		assert.Assert(t, rootDir.Has(tspath.Path("/src")), "/ should contain /src")
	})

	t.Run("deleted overlay removes from directory structure", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{}, false /* useCaseSensitiveFileNames */)

		// Previous overlay exists, now it's gone
		prevOverlays := map[tspath.Path]*Overlay{
			tspath.Path("/src/old.ts"): {
				fileBase: fileBase{fileName: "/src/old.ts", content: "const old = 1;"},
			},
		}
		overlays := make(map[tspath.Path]*Overlay)

		existingDirs := map[tspath.Path]*collections.Set[tspath.Path]{
			tspath.Path("/"):    collections.NewSetFromItems(tspath.Path("/src")),
			tspath.Path("/src"): collections.NewSetFromItems(tspath.Path("/src/old.ts")),
		}

		builder := newSnapshotFSBuilder(
			testFS,
			prevOverlays,
			overlays,
			make(map[tspath.Path]*diskFile),
			existingDirs,
			lsproto.PositionEncodingKindUTF16,
			toPath,
		)

		snapshot, _ := builder.Finalize()

		// Check that overlay file is removed from directory structure
		_, hasSrcDir := snapshot.directories[tspath.Path("/src")]
		assert.Assert(t, !hasSrcDir, "/src directory should be removed")

		_, hasRootDir := snapshot.directories[tspath.Path("/")]
		assert.Assert(t, !hasRootDir, "/ directory should be removed")
	})

	t.Run("overlay and disk file changes combined", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/lib/disk.ts": "const disk = 1;",
		}, false /* useCaseSensitiveFileNames */)

		// Add new overlay, remove old overlay, add disk file
		prevOverlays := map[tspath.Path]*Overlay{
			tspath.Path("/src/old.ts"): {
				fileBase: fileBase{fileName: "/src/old.ts", content: "const old = 1;"},
			},
		}
		overlays := map[tspath.Path]*Overlay{
			tspath.Path("/src/new.ts"): {
				fileBase: fileBase{fileName: "/src/new.ts", content: "const new = 1;"},
			},
		}

		existingDirs := map[tspath.Path]*collections.Set[tspath.Path]{
			tspath.Path("/"):    collections.NewSetFromItems(tspath.Path("/src")),
			tspath.Path("/src"): collections.NewSetFromItems(tspath.Path("/src/old.ts")),
		}

		builder := newSnapshotFSBuilder(
			testFS,
			prevOverlays,
			overlays,
			make(map[tspath.Path]*diskFile),
			existingDirs,
			lsproto.PositionEncodingKindUTF16,
			toPath,
		)

		// Add a disk file
		fh := builder.GetFile("/lib/disk.ts")
		assert.Assert(t, fh != nil)

		snapshot, changed := builder.Finalize()
		assert.Assert(t, changed, "should have changed")

		// /src should have new.ts but not old.ts
		srcDir, hasSrc := snapshot.directories[tspath.Path("/src")]
		assert.Assert(t, hasSrc, "/src directory should exist")
		assert.Assert(t, !srcDir.Has(tspath.Path("/src/old.ts")), "/src should not contain old.ts")
		assert.Assert(t, srcDir.Has(tspath.Path("/src/new.ts")), "/src should contain new.ts")

		// /lib should have disk.ts
		libDir, hasLib := snapshot.directories[tspath.Path("/lib")]
		assert.Assert(t, hasLib, "/lib directory should exist")
		assert.Assert(t, libDir.Has(tspath.Path("/lib/disk.ts")), "/lib should contain disk.ts")

		// root should have both /src and /lib
		rootDir := snapshot.directories[tspath.Path("/")]
		assert.Assert(t, rootDir.Has(tspath.Path("/src")), "root should contain /src")
		assert.Assert(t, rootDir.Has(tspath.Path("/lib")), "root should contain /lib")
	})

	t.Run("overlay in nested directory creates full path", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{}, false /* useCaseSensitiveFileNames */)

		prevOverlays := make(map[tspath.Path]*Overlay)
		overlays := map[tspath.Path]*Overlay{
			tspath.Path("/deep/nested/path/file.ts"): {
				fileBase: fileBase{fileName: "/deep/nested/path/file.ts", content: "const x = 1;"},
			},
		}

		builder := newSnapshotFSBuilder(
			testFS,
			prevOverlays,
			overlays,
			make(map[tspath.Path]*diskFile),
			make(map[tspath.Path]*collections.Set[tspath.Path]),
			lsproto.PositionEncodingKindUTF16,
			toPath,
		)

		snapshot, _ := builder.Finalize()

		// Check the complete directory tree was created
		assert.Assert(t, snapshot.directories[tspath.Path("/deep/nested/path")].Has(tspath.Path("/deep/nested/path/file.ts")))
		assert.Assert(t, snapshot.directories[tspath.Path("/deep/nested")].Has(tspath.Path("/deep/nested/path")))
		assert.Assert(t, snapshot.directories[tspath.Path("/deep")].Has(tspath.Path("/deep/nested")))
		assert.Assert(t, snapshot.directories[tspath.Path("/")].Has(tspath.Path("/deep")))
	})
}

func TestSnapshotFS(t *testing.T) {
	t.Parallel()

	toPath := func(fileName string) tspath.Path {
		return tspath.Path(fileName)
	}

	t.Run("GetFile returns overlay file", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/foo.ts": "disk content",
		}, false /* useCaseSensitiveFileNames */)

		overlays := map[tspath.Path]*Overlay{
			tspath.Path("/src/foo.ts"): {
				fileBase: fileBase{fileName: "/src/foo.ts", content: "overlay content"},
			},
		}

		snapshot := &SnapshotFS{
			toPath:      toPath,
			fs:          testFS,
			overlays:    overlays,
			diskFiles:   make(map[tspath.Path]*diskFile),
			directories: make(map[tspath.Path]*collections.Set[tspath.Path]),
		}

		fh := snapshot.GetFile("/src/foo.ts")
		assert.Assert(t, fh != nil)
		assert.Equal(t, fh.Content(), "overlay content")
	})

	t.Run("GetFile returns disk file when not in overlay", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/foo.ts": "disk content",
		}, false /* useCaseSensitiveFileNames */)

		diskFiles := map[tspath.Path]*diskFile{
			tspath.Path("/src/foo.ts"): newDiskFile("/src/foo.ts", "disk content"),
		}

		snapshot := &SnapshotFS{
			toPath:      toPath,
			fs:          testFS,
			overlays:    make(map[tspath.Path]*Overlay),
			diskFiles:   diskFiles,
			directories: make(map[tspath.Path]*collections.Set[tspath.Path]),
		}

		fh := snapshot.GetFile("/src/foo.ts")
		assert.Assert(t, fh != nil)
		assert.Equal(t, fh.Content(), "disk content")
	})

	t.Run("GetFile reads from fs when not cached", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/foo.ts": "fs content",
		}, false /* useCaseSensitiveFileNames */)

		snapshot := &SnapshotFS{
			toPath:      toPath,
			fs:          testFS,
			overlays:    make(map[tspath.Path]*Overlay),
			diskFiles:   make(map[tspath.Path]*diskFile),
			directories: make(map[tspath.Path]*collections.Set[tspath.Path]),
		}

		fh := snapshot.GetFile("/src/foo.ts")
		assert.Assert(t, fh != nil)
		assert.Equal(t, fh.Content(), "fs content")
	})

	t.Run("GetFile returns nil for non-existent file", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{}, false /* useCaseSensitiveFileNames */)

		snapshot := &SnapshotFS{
			toPath:      toPath,
			fs:          testFS,
			overlays:    make(map[tspath.Path]*Overlay),
			diskFiles:   make(map[tspath.Path]*diskFile),
			directories: make(map[tspath.Path]*collections.Set[tspath.Path]),
		}

		fh := snapshot.GetFile("/src/nonexistent.ts")
		assert.Assert(t, fh == nil, "should return nil for non-existent file")
	})

	t.Run("isOpenFile returns true for overlays", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{}, false /* useCaseSensitiveFileNames */)

		overlays := map[tspath.Path]*Overlay{
			tspath.Path("/src/foo.ts"): {
				fileBase: fileBase{fileName: "/src/foo.ts", content: "overlay content"},
			},
		}

		snapshot := &SnapshotFS{
			toPath:      toPath,
			fs:          testFS,
			overlays:    overlays,
			diskFiles:   make(map[tspath.Path]*diskFile),
			directories: make(map[tspath.Path]*collections.Set[tspath.Path]),
		}

		assert.Assert(t, snapshot.isOpenFile("/src/foo.ts"), "overlay file should be open")
		assert.Assert(t, !snapshot.isOpenFile("/src/bar.ts"), "non-overlay file should not be open")
	})

	t.Run("GetFileByPath uses provided path", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/foo.ts": "disk content",
		}, false /* useCaseSensitiveFileNames */)

		overlays := map[tspath.Path]*Overlay{
			tspath.Path("/src/foo.ts"): {
				fileBase: fileBase{fileName: "/src/foo.ts", content: "overlay content"},
			},
		}

		snapshot := &SnapshotFS{
			toPath:      toPath,
			fs:          testFS,
			overlays:    overlays,
			diskFiles:   make(map[tspath.Path]*diskFile),
			directories: make(map[tspath.Path]*collections.Set[tspath.Path]),
		}

		// GetFileByPath should use the provided path directly
		fh := snapshot.GetFileByPath("/src/foo.ts", tspath.Path("/src/foo.ts"))
		assert.Assert(t, fh != nil)
		assert.Equal(t, fh.Content(), "overlay content")
	})
}

func TestSourceFS(t *testing.T) {
	t.Parallel()

	toPath := func(fileName string) tspath.Path {
		return tspath.Path(fileName)
	}

	t.Run("tracks files when tracking enabled", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/foo.ts": "content",
		}, false /* useCaseSensitiveFileNames */)

		snapshot := &SnapshotFS{
			toPath:      toPath,
			fs:          testFS,
			overlays:    make(map[tspath.Path]*Overlay),
			diskFiles:   make(map[tspath.Path]*diskFile),
			directories: make(map[tspath.Path]*collections.Set[tspath.Path]),
		}

		sourceFS := newSourceFS(true /* tracking */, snapshot, toPath)

		// File should not be seen yet
		assert.Assert(t, !sourceFS.Seen(tspath.Path("/src/foo.ts")))

		// Read the file
		fh := sourceFS.GetFile("/src/foo.ts")
		assert.Assert(t, fh != nil)

		// Now it should be seen
		assert.Assert(t, sourceFS.Seen(tspath.Path("/src/foo.ts")))
	})

	t.Run("does not track files when tracking disabled", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/foo.ts": "content",
		}, false /* useCaseSensitiveFileNames */)

		snapshot := &SnapshotFS{
			toPath:      toPath,
			fs:          testFS,
			overlays:    make(map[tspath.Path]*Overlay),
			diskFiles:   make(map[tspath.Path]*diskFile),
			directories: make(map[tspath.Path]*collections.Set[tspath.Path]),
		}

		sourceFS := newSourceFS(false /* tracking */, snapshot, toPath)

		// Read the file
		fh := sourceFS.GetFile("/src/foo.ts")
		assert.Assert(t, fh != nil)

		// Should not be seen since tracking is disabled
		assert.Assert(t, !sourceFS.Seen(tspath.Path("/src/foo.ts")))
	})

	t.Run("DisableTracking stops tracking", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/foo.ts": "content",
			"/src/bar.ts": "content",
		}, false /* useCaseSensitiveFileNames */)

		snapshot := &SnapshotFS{
			toPath:      toPath,
			fs:          testFS,
			overlays:    make(map[tspath.Path]*Overlay),
			diskFiles:   make(map[tspath.Path]*diskFile),
			directories: make(map[tspath.Path]*collections.Set[tspath.Path]),
		}

		sourceFS := newSourceFS(true /* tracking */, snapshot, toPath)

		// Read foo while tracking
		sourceFS.GetFile("/src/foo.ts")
		assert.Assert(t, sourceFS.Seen(tspath.Path("/src/foo.ts")))

		// Disable tracking
		sourceFS.DisableTracking()

		// Read bar after tracking disabled
		sourceFS.GetFile("/src/bar.ts")
		assert.Assert(t, !sourceFS.Seen(tspath.Path("/src/bar.ts")))
	})

	t.Run("FileExists returns true for files in source", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/foo.ts": "content",
		}, false /* useCaseSensitiveFileNames */)

		snapshot := &SnapshotFS{
			toPath:      toPath,
			fs:          testFS,
			overlays:    make(map[tspath.Path]*Overlay),
			diskFiles:   make(map[tspath.Path]*diskFile),
			directories: make(map[tspath.Path]*collections.Set[tspath.Path]),
		}

		sourceFS := newSourceFS(false /* tracking */, snapshot, toPath)

		assert.Assert(t, sourceFS.FileExists("/src/foo.ts"))
		assert.Assert(t, !sourceFS.FileExists("/src/nonexistent.ts"))
	})

	t.Run("ReadFile returns content for files in source", func(t *testing.T) {
		t.Parallel()
		testFS := vfstest.FromMap(map[string]string{
			"/src/foo.ts": "file content",
		}, false /* useCaseSensitiveFileNames */)

		snapshot := &SnapshotFS{
			toPath:      toPath,
			fs:          testFS,
			overlays:    make(map[tspath.Path]*Overlay),
			diskFiles:   make(map[tspath.Path]*diskFile),
			directories: make(map[tspath.Path]*collections.Set[tspath.Path]),
		}

		sourceFS := newSourceFS(false /* tracking */, snapshot, toPath)

		content, ok := sourceFS.ReadFile("/src/foo.ts")
		assert.Assert(t, ok)
		assert.Equal(t, content, "file content")

		_, ok = sourceFS.ReadFile("/src/nonexistent.ts")
		assert.Assert(t, !ok)
	})
}
