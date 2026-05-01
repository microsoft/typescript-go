package vfswatch_test

import (
	"slices"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"github.com/microsoft/typescript-go/internal/vfs/vfswatch"
)

func TestWatchEventPipeline(t *testing.T) {
	t.Parallel()
	files := map[string]string{
		"/src/a.ts":      "const a = 1;",
		"/src/b.ts":      "const b = 2;",
		"/tsconfig.json": "{}",
	}
	fs := vfstest.FromMap(files, true)

	fw := vfswatch.NewFileWatcher(fs, 50*time.Millisecond, true, func(e vfswatch.WatchEvent) {})
	fw.UpdateWatchedDirectories(map[string]bool{"/src": true})

	// No changes after initial snapshot
	event := fw.ScanForChanges()
	if event.HasChanges() {
		t.Fatalf("unexpected changes: created=%v deleted=%v changed=%v", event.Created, event.Deleted, event.Changed)
	}

	// Sleep to ensure WriteFile gets a different modTime than the snapshot,
	// since Windows time resolution can be ~15ms.
	time.Sleep(20 * time.Millisecond)

	// Modify a file → detect change
	_ = fs.WriteFile("/src/a.ts", "const a = 999;")
	event = fw.ScanForChanges()
	if !event.HasChanges() {
		t.Fatal("no changes detected after modifying file")
	}
	if len(event.Changed) == 0 {
		t.Fatalf("expected Changed entries, got: created=%v deleted=%v changed=%v", event.Created, event.Deleted, event.Changed)
	}
	if !slices.Contains(event.Changed, "/src/a.ts") {
		t.Fatalf("expected /src/a.ts in Changed, got: %v", event.Changed)
	}
	fw.UpdateWatchedDirectories(map[string]bool{"/src": true})
	time.Sleep(20 * time.Millisecond)

	// Add new file → detect create
	_ = fs.WriteFile("/src/new.ts", "const n = 42;")
	event = fw.ScanForChanges()
	if !event.HasChanges() {
		t.Fatal("no changes detected after adding file")
	}
	if len(event.Created) == 0 {
		t.Fatalf("expected Created entries, got: created=%v deleted=%v changed=%v", event.Created, event.Deleted, event.Changed)
	}
	fw.UpdateWatchedDirectories(map[string]bool{"/src": true})
	time.Sleep(20 * time.Millisecond)

	// Delete a file → detect delete
	_ = fs.Remove("/src/b.ts")
	event = fw.ScanForChanges()
	if !event.HasChanges() {
		t.Fatal("no changes detected after deleting file")
	}
	if len(event.Deleted) == 0 {
		t.Fatalf("expected Deleted entries, got: created=%v deleted=%v changed=%v", event.Created, event.Deleted, event.Changed)
	}
	fw.UpdateWatchedDirectories(map[string]bool{"/src": true})

	// No changes after refresh
	event = fw.ScanForChanges()
	if event.HasChanges() {
		t.Fatalf("unexpected changes after refresh: created=%v deleted=%v changed=%v", event.Created, event.Deleted, event.Changed)
	}
	time.Sleep(20 * time.Millisecond)

	// Add subdirectory with file → detect create
	_ = fs.WriteFile("/src/sub/deep.ts", "const deep = true;")
	event = fw.ScanForChanges()
	if !event.HasChanges() {
		t.Fatal("no changes detected after adding file in subdirectory")
	}
	if len(event.Created) == 0 {
		t.Fatalf("expected Created entries, got: created=%v deleted=%v changed=%v", event.Created, event.Deleted, event.Changed)
	}
}

func TestWatchEventNonRecursive(t *testing.T) {
	t.Parallel()
	files := map[string]string{
		"/root/a.ts":     "const a = 1;",
		"/root/sub/b.ts": "const b = 2;",
	}
	fs := vfstest.FromMap(files, true)

	fw := vfswatch.NewFileWatcher(fs, 50*time.Millisecond, true, func(e vfswatch.WatchEvent) {})
	fw.UpdateWatchedDirectories(map[string]bool{"/root": false}) // non-recursive

	// Sleep to ensure WriteFile gets a different modTime than the snapshot,
	// since Windows time resolution can be ~15ms.
	time.Sleep(20 * time.Millisecond)

	// Modify a file in root → detect change
	_ = fs.WriteFile("/root/a.ts", "const a = 999;")
	event := fw.ScanForChanges()
	if !event.HasChanges() {
		t.Fatal("no changes detected after modifying file in non-recursive dir")
	}
	fw.UpdateWatchedDirectories(map[string]bool{"/root": false})
	time.Sleep(20 * time.Millisecond)

	// Add file in root → detect create
	_ = fs.WriteFile("/root/c.ts", "const c = 3;")
	event = fw.ScanForChanges()
	if !event.HasChanges() {
		t.Fatal("no changes detected after adding file in non-recursive dir")
	}
	if len(event.Created) == 0 {
		t.Fatalf("expected Created entries, got: created=%v deleted=%v changed=%v", event.Created, event.Deleted, event.Changed)
	}
}
