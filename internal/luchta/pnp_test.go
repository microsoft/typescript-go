package luchta

import (
	"os"
	"path/filepath"
	"testing"
)

// pnpFixtureCjs locates the simplest PnP fixture shipped with the repo.
func pnpFixtureCjs(t *testing.T) string {
	t.Helper()
	// repo-relative: testdata/fixtures/pnp/pnp-yarn-v4.cjs
	root, err := filepath.Abs(filepath.Join("..", "..", "testdata", "fixtures", "pnp"))
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(root, "pnp-yarn-v4.cjs")
	if !fileExists(p) {
		t.Skipf("PnP fixture not found at %s", p)
	}
	return p
}

func TestCompilerFSNoPnp(t *testing.T) {
	cwd := t.TempDir() // no .pnp.cjs anywhere above a fresh temp dir
	fsys, api := compilerFS(cwd)
	if api != nil {
		t.Fatalf("expected nil PnP API outside a PnP workspace")
	}
	if fsys == nil {
		t.Fatalf("expected a non-nil base FS")
	}
}

func TestCompilerFSDetectsPnp(t *testing.T) {
	root := t.TempDir()
	// Copy a real fixture manifest to the workspace root.
	// InitPnpApi requires only that .pnp.cjs is present and parseable (it reads
	// and parses the JSON payload embedded in the CJS file; it does NOT require
	// that any referenced package directories exist). A copied fixture is sufficient.
	data, err := os.ReadFile(pnpFixtureCjs(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".pnp.cjs"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	pkg := filepath.Join(root, "packages", "app")
	if err := os.MkdirAll(pkg, 0o755); err != nil {
		t.Fatal(err)
	}
	_, api := compilerFS(pkg)
	if api == nil {
		t.Fatalf("expected non-nil PnP API when .pnp.cjs is present above cwd")
	}
}
