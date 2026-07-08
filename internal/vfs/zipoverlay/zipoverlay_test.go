package zipoverlay

import (
	"archive/zip"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

// writeFixtureZip creates a zip archive that looks like a Yarn cache entry:
// node_modules/pkg/{package.json,index.js,cjs/impl.js}.
func writeFixtureZip(t *testing.T, path string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	w := zip.NewWriter(f)
	entries := map[string]string{
		"node_modules/pkg/package.json": `{"name":"pkg","version":"1.0.0"}`,
		"node_modules/pkg/index.js":     "module.exports = 1;\n",
		"node_modules/pkg/cjs/impl.js":  "module.exports = 2;\n",
	}
	for name, body := range entries {
		fw, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := fw.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestZipOverlay(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	zipPath := filepath.ToSlash(filepath.Join(dir, "pkg-npm-1.0.0-abc.zip"))
	writeFixtureZip(t, zipPath)

	o := Wrap(osvfs.FS())
	pkgDir := zipPath + "/node_modules/pkg"

	// The .zip file and its interior directories all report as directories, so
	// the module resolver's parent-directory checks pass.
	for _, d := range []string{zipPath, zipPath + "/node_modules", pkgDir, pkgDir + "/cjs"} {
		if !o.DirectoryExists(d) {
			t.Errorf("DirectoryExists(%q) = false, want true", d)
		}
	}

	// Files inside the zip exist and are readable; a file is not a directory.
	indexJS := pkgDir + "/index.js"
	if !o.FileExists(indexJS) {
		t.Errorf("FileExists(%q) = false", indexJS)
	}
	if o.DirectoryExists(indexJS) {
		t.Errorf("a file reported as a directory: %q", indexJS)
	}
	if c, ok := o.ReadFile(pkgDir + "/package.json"); !ok || c != `{"name":"pkg","version":"1.0.0"}` {
		t.Errorf("ReadFile package.json = %q, ok=%v", c, ok)
	}
	if fi := o.Stat(indexJS); fi == nil || fi.IsDir() || fi.Size() == 0 {
		t.Errorf("Stat(index.js) wrong: %+v", fi)
	}

	// A missing file inside the zip does not exist.
	if o.FileExists(pkgDir + "/missing.js") {
		t.Error("missing file reported as existing")
	}

	// GetAccessibleEntries lists immediate children only.
	entries := o.GetAccessibleEntries(pkgDir)
	if !slices.Contains(entries.Files, "index.js") || !slices.Contains(entries.Files, "package.json") {
		t.Errorf("entries.Files missing expected: %v", entries.Files)
	}
	if !slices.Contains(entries.Directories, "cjs") {
		t.Errorf("entries.Directories missing cjs: %v", entries.Directories)
	}
	if slices.Contains(entries.Files, "impl.js") {
		t.Error("GetAccessibleEntries returned a non-immediate child")
	}

	// WalkDir visits every entry inside the zip.
	seen := map[string]bool{}
	if err := o.WalkDir(pkgDir, func(path string, d os.DirEntry, err error) error {
		seen[path] = true
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if !seen[indexJS] || !seen[pkgDir+"/cjs/impl.js"] {
		t.Errorf("WalkDir did not visit all entries: %v", seen)
	}

	// A non-zip path delegates to the inner FS: the .zip file itself is a real
	// file on disk, and a missing sibling reads back as not-ok.
	if !osvfs.FS().FileExists(zipPath) {
		t.Fatal("fixture zip missing on disk")
	}
	if _, ok := o.ReadFile(dir + "/does-not-exist"); ok {
		t.Error("delegated ReadFile of a missing file returned ok")
	}
}

// TestNonZipVirtualDelegates covers a PnP __virtual__ path whose backing
// location is a plain on-disk directory (a workspace or unplugged package with a
// peer dependency), not a cache .zip. The overlay must dereference the virtual
// segment and delegate the resulting real path to the inner FS. Regression for
// the bug where split() fast-rejected any path without ".zip" before
// dereferencing, so these paths returned "not found" and the resolver reported
// TS2307 for a workspace library that has a peer dependency.
func TestNonZipVirtualDelegates(t *testing.T) {
	t.Parallel()
	dir := filepath.ToSlash(t.TempDir())
	libDir := filepath.Join(dir, "lib")
	if err := os.MkdirAll(libDir, 0o755); err != nil {
		t.Fatal(err)
	}
	realFile := filepath.Join(libDir, "index.d.ts")
	if err := os.WriteFile(realFile, []byte("export const ui: number;\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	o := Wrap(osvfs.FS())
	// <dir>/x/__virtual__/<hash>/1/lib/index.d.ts dereferences to <dir>/lib/index.d.ts
	// (n=1 pops the "x" segment before the __virtual__ triple).
	virtual := dir + "/x/__virtual__/w-ui-virtual-aaa/1/lib/index.d.ts"
	if !o.FileExists(virtual) {
		t.Fatalf("FileExists(%q) = false; virtual path did not dereference to the real file", virtual)
	}
	if c, ok := o.ReadFile(virtual); !ok || c != "export const ui: number;\n" {
		t.Fatalf("ReadFile(virtual) = %q ok=%v", c, ok)
	}
	// The virtual directory dereferences too.
	virtualDir := dir + "/x/__virtual__/w-ui-virtual-aaa/1/lib"
	if !o.DirectoryExists(virtualDir) {
		t.Errorf("DirectoryExists(%q) = false", virtualDir)
	}
	// Realpath keeps the virtual path in its own space (no desync with the PnP
	// locator table).
	if got := o.Realpath(virtual); got != virtual {
		t.Errorf("Realpath(virtual) = %q, want identity", got)
	}
}

// TestZipBackedVirtualReadsFromArchive covers a PnP __virtual__ path whose
// backing location is inside a cache .zip (an npm package with a peer
// dependency). The overlay dereferences the virtual segment to the .zip path and
// serves the file from the archive, while Realpath keeps the virtual path.
func TestZipBackedVirtualReadsFromArchive(t *testing.T) {
	t.Parallel()
	dir := filepath.ToSlash(t.TempDir())
	cacheDir := filepath.Join(dir, ".yarn", "cache")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		t.Fatal(err)
	}
	zipPath := filepath.ToSlash(filepath.Join(cacheDir, "pkg-npm-1.0.0-abc.zip"))
	writeFixtureZip(t, zipPath)

	o := Wrap(osvfs.FS())
	// The Yarn shape for a peer-virtualized npm package: the cache .zip lives at
	// .yarn/cache/, and the virtual location is
	// .yarn/__virtual__/<hash>/0/cache/<zip>/node_modules/pkg/. Dereferencing
	// (count 0) pops nothing and rejoins onto .yarn/, landing on the real .zip.
	virtual := dir + "/.yarn/__virtual__/pkg-virtual-abc/0/cache/pkg-npm-1.0.0-abc.zip/node_modules/pkg/index.js"
	if !o.FileExists(virtual) {
		t.Fatalf("FileExists(%q) = false; zip-backed virtual path not served", virtual)
	}
	if c, ok := o.ReadFile(virtual); !ok || c != "module.exports = 1;\n" {
		t.Fatalf("ReadFile(zip-backed virtual) = %q ok=%v", c, ok)
	}
	if got := o.Realpath(virtual); got != virtual {
		t.Errorf("Realpath(virtual) = %q, want identity", got)
	}
}

// TestZipEdgeCases covers the overlay's boundary conditions: a path referencing a
// .zip that does not exist on disk falls through to the inner FS (not a crash or a
// false hit), the legacy $$virtual spelling is handled, and a path that merely
// contains ".zip" as a substring is not treated as an archive.
func TestZipEdgeCases(t *testing.T) {
	t.Parallel()
	dir := filepath.ToSlash(t.TempDir())
	o := Wrap(osvfs.FS())

	// A .zip archive that does not exist on disk: split() finds no real archive and
	// delegates, so the interior path simply does not exist (no panic).
	ghost := dir + "/missing-npm-1.0.0.zip/node_modules/pkg/index.js"
	if o.FileExists(ghost) {
		t.Error("interior path of a non-existent .zip reported as existing")
	}
	if _, ok := o.ReadFile(ghost); ok {
		t.Error("ReadFile of a non-existent .zip interior returned ok")
	}

	// A real file whose name merely contains ".zip" (not "<name>.zip/…" and not a
	// ".zip" suffix) is delegated, read as an ordinary file.
	plain := filepath.Join(dir, "notes.zip.txt")
	if err := os.WriteFile(plain, []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if c, ok := o.ReadFile(filepath.ToSlash(plain)); !ok || c != "hello\n" {
		t.Errorf("delegated read of a .zip-substring file = %q ok=%v", c, ok)
	}

	// The legacy $$virtual spelling dereferences like __virtual__.
	libDir := filepath.Join(dir, "lib")
	if err := os.MkdirAll(libDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libDir, "i.d.ts"), []byte("export const z=1;\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	legacy := dir + "/x/$$virtual/h/1/lib/i.d.ts"
	if !o.FileExists(legacy) {
		t.Errorf("FileExists(%q) = false; $$virtual not dereferenced", legacy)
	}
}
