package luchta

import (
	"bytes"
	"context"
	"path/filepath"
	"slices"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

func writeTsPackage(t *testing.T, dir, tsconfig, srcName, srcBody string) {
	t.Helper()
	mustWrite(t, filepath.Join(dir, "tsconfig.json"), tsconfig)
	mustWrite(t, filepath.Join(dir, "src", srcName), srcBody)
}

func TestCompilePackageSuccess(t *testing.T) {
	cwd := t.TempDir()
	writeTsPackage(t, cwd, `{
		"compilerOptions": {"declaration": true, "outDir": "dist", "rootDir": "src", "module": "nodenext", "moduleResolution": "nodenext"},
		"include": ["src/**/*"]
	}`, "index.ts", "export const answer: number = 42;\n")

	res := CompilePackage(context.Background(), cwd)
	if res.ExitCode != 0 {
		t.Fatalf("exitCode=%d diagnostics=%s", res.ExitCode, res.Diagnostics)
	}
	if !fileExists(filepath.Join(cwd, "dist", "index.js")) {
		t.Fatalf("expected dist/index.js emitted; outputs=%v", res.Outputs)
	}
	if !slices.Contains(res.Outputs, "dist/index.js") {
		t.Fatalf("outputs missing dist/index.js: %v", res.Outputs)
	}
	if !slices.Contains(res.Inputs, "tsconfig.json") {
		t.Fatalf("inputs missing tsconfig.json: %v", res.Inputs)
	}
}

func TestCompilePackageTypeError(t *testing.T) {
	cwd := t.TempDir()
	writeTsPackage(t, cwd, `{
		"compilerOptions": {"outDir": "dist", "rootDir": "src", "module": "nodenext", "moduleResolution": "nodenext", "strict": true}
	}`, "index.ts", "export const answer: number = \"not a number\";\n")

	res := CompilePackage(context.Background(), cwd)
	if res.ExitCode == 0 {
		t.Fatalf("expected non-zero exit on type error")
	}
	if res.Diagnostics == "" {
		t.Fatalf("expected diagnostic text")
	}
}

func TestCompilePackageNoTsconfig(t *testing.T) {
	cwd := t.TempDir()
	res := CompilePackage(context.Background(), cwd)
	if res.ExitCode != 0 {
		t.Fatalf("missing tsconfig should be a no-op success, got %d", res.ExitCode)
	}
	if !slices.Contains(res.Inputs, "src/**") {
		t.Fatalf("expected default src/** input, got %v", res.Inputs)
	}
}

// TestCompilePackageConcurrentRace verifies that CompilePackage is race-free when
// called concurrently. Run with `go test -race` to exercise the parallel emit
// and the goroutine-safe EmittedFiles assembly.
func TestCompilePackageConcurrentRace(t *testing.T) {
	const goroutines = 4

	// Build a package with three source files so the emitter spawns multiple
	// parallel goroutines internally.
	makePkg := func(t *testing.T) string {
		t.Helper()
		cwd := t.TempDir()
		tsconfig := `{
			"compilerOptions": {"declaration": true, "outDir": "dist", "rootDir": "src", "module": "nodenext", "moduleResolution": "nodenext"},
			"include": ["src/**/*"]
		}`
		mustWrite(t, filepath.Join(cwd, "tsconfig.json"), tsconfig)
		mustWrite(t, filepath.Join(cwd, "src", "a.ts"), "export const a = 1;\n")
		mustWrite(t, filepath.Join(cwd, "src", "b.ts"), "export const b = 2;\n")
		mustWrite(t, filepath.Join(cwd, "src", "c.ts"), "export const c = 3;\n")
		return cwd
	}

	// Each goroutine gets its own isolated package directory so they don't
	// race on the filesystem either.
	cwds := make([]string, goroutines)
	for i := range cwds {
		cwds[i] = makePkg(t)
	}

	type result struct {
		res CompileResult
		idx int
	}
	ch := make(chan result, goroutines)

	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			ch <- result{CompilePackage(context.Background(), cwds[i]), i}
		}()
	}

	for range goroutines {
		r := <-ch
		if r.res.ExitCode != 0 {
			t.Errorf("goroutine %d: exitCode=%d diagnostics=%s", r.idx, r.res.ExitCode, r.res.Diagnostics)
		}
		if len(r.res.Outputs) == 0 {
			t.Errorf("goroutine %d: expected non-empty Outputs", r.idx)
		}
	}
}

// TestReportDiagnosticsNilParsedNoPanic verifies that reportDiagnostics does not
// panic when parsed is nil (e.g. GetParsedCommandLineOfConfigFile returns nil on a
// fatal config-read error). This guards against the nil-pointer dereference in
// CreateDiagnosticReporter which immediately accesses options.Quiet.
func TestReportDiagnosticsNilParsedNoPanic(t *testing.T) {
	cwd := t.TempDir()
	fsys := bundled.WrapFS(osvfs.FS())
	var buf bytes.Buffer
	sys := newRunSystem(cwd, fsys, bundled.LibPath(), &buf, nil)

	// Must not panic even with nil parsed and an empty diagnostics slice.
	reportDiagnostics(sys, &buf, nil, []*ast.Diagnostic{})
}
