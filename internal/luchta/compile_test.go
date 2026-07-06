package luchta

import (
	"context"
	"encoding/json"
	"path/filepath"
	"slices"
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
		t.Fatalf("invalid SARIF: %v", err)
	}
	if len(log.Runs) == 0 || len(log.Runs[0].Results) == 0 {
		t.Fatalf("SARIF missing results")
	}
	result := log.Runs[0].Results[0]
	if result.RuleID == "" {
		t.Fatalf("missing ruleId")
	}
	if len(result.Locations) == 0 {
		t.Fatalf("missing location")
	}
	loc := result.Locations[0]
	if loc.PhysicalLocation.ArtifactLocation.URI == "" {
		t.Fatalf("missing artifact URI")
	}
}

func TestCompilePackageWithTsconfigBuild(t *testing.T) {
	cwd := t.TempDir()
	tsconfigContent := `{
		"compilerOptions": {"declaration": true, "outDir": "dist", "rootDir": "src", "module": "nodenext", "moduleResolution": "nodenext"},
		"include": ["src/**/*"]
	}`
	mustWrite(t, filepath.Join(cwd, "tsconfig.json"), `{}`)
	mustWrite(t, filepath.Join(cwd, "tsconfig.build.json"), tsconfigContent)
	mustWrite(t, filepath.Join(cwd, "src", "index.ts"), "export const x = 1;\n")

	res := CompilePackage(context.Background(), cwd)
	if res.ExitCode != 0 {
		t.Fatalf("exitCode=%d", res.ExitCode)
	}
	if !slices.Contains(res.Outputs, "dist/index.js") {
		t.Fatalf("outputs missing dist/index.js: %v", res.Outputs)
	}
}

func TestResolveInputs(t *testing.T) {
	cwd := t.TempDir()
	writeTsPackage(t, cwd, `{
		"compilerOptions": {"outDir": "dist"},
		"include": ["src/**"]
	}`, "index.ts", "export const x = 1;\n")

	inputs, err := ResolveInputs(cwd)
	if err != nil {
		t.Fatalf("ResolveInputs: %v", err)
	}
	// Should contain tsconfig.json
	if !slices.Contains(inputs, "tsconfig.json") {
		t.Fatalf("inputs missing tsconfig.json: %v", inputs)
	}
	// Should contain tsconfig.build.json candidate
	if !slices.Contains(inputs, "tsconfig.build.json") {
		t.Fatalf("inputs missing tsconfig.build.json candidate: %v", inputs)
	}
	// Should contain include pattern
	if !slices.Contains(inputs, "src/**") {
		t.Fatalf("inputs missing src/** pattern: %v", inputs)
	}
}

func TestResolveInputsNoTsconfig(t *testing.T) {
	cwd := t.TempDir()
	// No tsconfig files exist
	inputs, err := ResolveInputs(cwd)
	if err != nil {
		t.Fatalf("ResolveInputs: %v", err)
	}
	// Should still report candidates as absent literals
	for _, name := range tsconfigCandidates {
		if !slices.Contains(inputs, name) {
			t.Fatalf("inputs missing absent candidate %q: %v", name, inputs)
		}
	}
	// When no tsconfig exists, includePatterns returns default "src/**"
	// which gets added, so we should see it
	if !slices.Contains(inputs, "src/**") {
		t.Fatalf("inputs missing default src/**: %v", inputs)
	}
}

func TestResolveInputsIncludeGlobWithoutMatchingFiles(t *testing.T) {
	cwd := t.TempDir()
	// tsconfig declares an include glob, but NO matching source files exist on
	// disk. ResolveInputs must still return the declared glob pattern, proving
	// it reads the raw include specs rather than enumerating files.
	mustWrite(t, filepath.Join(cwd, "tsconfig.json"), `{
		"compilerOptions": {"outDir": "dist"},
		"include": ["lib/**/*.ts", "types/**"]
	}`)

	inputs, err := ResolveInputs(cwd)
	if err != nil {
		t.Fatalf("ResolveInputs: %v", err)
	}
	for _, want := range []string{"lib/**/*.ts", "types/**", "tsconfig.json", "tsconfig.build.json"} {
		if !slices.Contains(inputs, want) {
			t.Fatalf("inputs missing %q: %v", want, inputs)
		}
	}
	// Must NOT contain the default fallback when include is explicitly set.
	if slices.Contains(inputs, "src/**") {
		t.Fatalf("unexpected default src/** when include is explicit: %v", inputs)
	}
}

func TestResolveInputsIncludeFromExtends(t *testing.T) {
	cwd := t.TempDir()
	// The child tsconfig has no include of its own; it inherits include from an
	// extended base config. ResolveInputs must resolve `extends` and surface the
	// inherited include pattern.
	mustWrite(t, filepath.Join(cwd, "tsconfig.base.json"), `{
		"include": ["packages/**/*.ts"]
	}`)
	mustWrite(t, filepath.Join(cwd, "tsconfig.json"), `{
		"extends": "./tsconfig.base.json",
		"compilerOptions": {"outDir": "dist"}
	}`)

	inputs, err := ResolveInputs(cwd)
	if err != nil {
		t.Fatalf("ResolveInputs: %v", err)
	}
	if !slices.Contains(inputs, "packages/**/*.ts") {
		t.Fatalf("inputs missing include inherited via extends: %v", inputs)
	}
	// The extended (base) config lives in the SAME package dir, so its path must
	// be reported as an input so edits to the base config bust the cache.
	if !slices.Contains(inputs, "tsconfig.base.json") {
		t.Fatalf("inputs missing same-package extended base config path: %v", inputs)
	}
}

func TestResolveInputsExtendsFromSubdirIsTracked(t *testing.T) {
	cwd := t.TempDir()
	// Base config in a subdirectory of the package (still same-package,
	// non-escaping) must be tracked as an input.
	mustWrite(t, filepath.Join(cwd, "config", "tsconfig.base.json"), `{
		"include": ["lib/**/*.ts"]
	}`)
	mustWrite(t, filepath.Join(cwd, "tsconfig.json"), `{
		"extends": "./config/tsconfig.base.json",
		"compilerOptions": {"outDir": "dist"}
	}`)

	inputs, err := ResolveInputs(cwd)
	if err != nil {
		t.Fatalf("ResolveInputs: %v", err)
	}
	// tsoptions rebases the base config's include patterns relative to the base
	// config's own directory, so "lib/**/*.ts" from config/ becomes
	// "config/lib/**/*.ts".
	if !slices.Contains(inputs, "config/lib/**/*.ts") {
		t.Fatalf("inputs missing rebased include inherited via subdir extends: %v", inputs)
	}
	if !slices.Contains(inputs, "config/tsconfig.base.json") {
		t.Fatalf("inputs missing subdir extended base config path: %v", inputs)
	}
}

func TestCompilePackageConcurrentRace(t *testing.T) {
	const goroutines = 4

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
