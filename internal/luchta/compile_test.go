package luchta

import (
	"context"
	"encoding/json"
	"path/filepath"
	"slices"
	"strings"
	"testing"
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
		t.Fatalf("exitCode=%d sarif=%s", res.ExitCode, DiagnosticsToSARIF(res.Diagnostics))
	}
	if len(res.Diagnostics) != 0 {
		t.Fatalf("expected no diagnostics on success, got %s", DiagnosticsToSARIF(res.Diagnostics))
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
	if len(res.Diagnostics) == 0 {
		t.Fatalf("expected diagnostics")
	}

	// The worker reports diagnostics as a SARIF document; verify it is
	// well-formed and carries the error with a clickable location.
	var log sarifLog
	if err := json.Unmarshal([]byte(DiagnosticsToSARIF(res.Diagnostics)), &log); err != nil {
		t.Fatalf("SARIF output is not valid JSON: %v", err)
	}
	if log.Version != "2.1.0" {
		t.Fatalf("unexpected SARIF version %q", log.Version)
	}
	if len(log.Runs) != 1 || len(log.Runs[0].Results) == 0 {
		t.Fatalf("expected at least one SARIF result, got %+v", log.Runs)
	}
	got := log.Runs[0].Results[0]
	if got.Level != "error" {
		t.Fatalf("expected error level, got %q", got.Level)
	}
	if !strings.HasPrefix(got.RuleID, "TS") {
		t.Fatalf("expected TS-prefixed ruleId, got %q", got.RuleID)
	}
	if got.Message.Text == "" {
		t.Fatalf("expected a non-empty message")
	}
	if len(got.Locations) == 0 || got.Locations[0].PhysicalLocation.Region == nil {
		t.Fatalf("expected a physical location with a region, got %+v", got.Locations)
	}
	if region := got.Locations[0].PhysicalLocation.Region; region.StartLine < 1 || region.StartColumn < 1 {
		t.Fatalf("expected 1-based region, got %+v", region)
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
			t.Errorf("goroutine %d: exitCode=%d sarif=%s", r.idx, r.res.ExitCode, DiagnosticsToSARIF(r.res.Diagnostics))
		}
		if len(r.res.Outputs) == 0 {
			t.Errorf("goroutine %d: expected non-empty Outputs", r.idx)
		}
	}
}
