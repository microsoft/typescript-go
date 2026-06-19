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
