package luchta

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestRelativizeOutputs(t *testing.T) {
	got := RelativizeOutputs("/repo/pkg", []string{"/repo/pkg/dist/a.js", "/repo/pkg/dist/types/a.d.ts"})
	want := []string{"dist/a.js", "dist/types/a.d.ts"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestCleanOutputsRemovesStaleDts(t *testing.T) {
	cwd := t.TempDir()
	mustWrite(t, filepath.Join(cwd, "src", "keep.ts"), "export const x = 1;")
	mustWrite(t, filepath.Join(cwd, "dist/types", "keep.d.ts"), "export declare const x: number;")
	mustWrite(t, filepath.Join(cwd, "dist/types", "keep.d.ts.map"), "{}")
	mustWrite(t, filepath.Join(cwd, "dist/types", "gone.d.ts"), "export declare const y: number;")
	mustWrite(t, filepath.Join(cwd, "dist/types", "gone.d.ts.map"), "{}")

	if err := CleanOutputs(cwd, "dist/types", "src", false); err != nil {
		t.Fatalf("CleanOutputs: %v", err)
	}
	got := listFiles(t, filepath.Join(cwd, "dist/types"))
	want := []string{"keep.d.ts", "keep.d.ts.map"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestCleanOutputsSkipsWhenNoEmit(t *testing.T) {
	cwd := t.TempDir()
	mustWrite(t, filepath.Join(cwd, "dist/types", "gone.d.ts"), "x")
	if err := CleanOutputs(cwd, "dist/types", "src", true); err != nil {
		t.Fatalf("CleanOutputs: %v", err)
	}
	if _, err := os.Stat(filepath.Join(cwd, "dist/types", "gone.d.ts")); err != nil {
		t.Fatalf("noEmit should leave files untouched: %v", err)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func listFiles(t *testing.T, dir string) []string {
	t.Helper()
	var out []string
	_ = filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			rel, _ := filepath.Rel(dir, p)
			out = append(out, rel)
		}
		return nil
	})
	sort.Strings(out)
	return out
}
