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
	clock := vfstest.NewSteppingClock(100 * time.Millisecond)
	fs := vfstest.FromMapWithClock(files, true, clock)

	fw := vfswatch.NewFileWatcher(fs, 50*time.Millisecond, true, func(e vfswatch.WatchEvent) {})
	fw.UpdateWatchedDirectories(map[string]bool{"/src": true})

	// No changes after initial snapshot
	event := fw.ScanForChanges()
	if event.HasChanges() {
		t.Fatalf("unexpected changes: created=%v deleted=%v changed=%v", event.Created, event.Deleted, event.Changed)
	}

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
	clock := vfstest.NewSteppingClock(100 * time.Millisecond)
	fs := vfstest.FromMapWithClock(files, true, clock)

	fw := vfswatch.NewFileWatcher(fs, 50*time.Millisecond, true, func(e vfswatch.WatchEvent) {})
	fw.UpdateWatchedDirectories(map[string]bool{"/root": false}) // non-recursive

	// Modify a file in root → detect change
	_ = fs.WriteFile("/root/a.ts", "const a = 999;")
	event := fw.ScanForChanges()
	if !event.HasChanges() {
		t.Fatal("no changes detected after modifying file in non-recursive dir")
	}
	fw.UpdateWatchedDirectories(map[string]bool{"/root": false})

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
