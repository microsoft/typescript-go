# luchta-tsc-worker Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a statically-linked Go binary, `luchta-tsc-worker`, that speaks luchta's JSONL-over-stdio worker protocol, compiles each package with tsgo's compiler internals, and resolves modules via Yarn PnP — cross-compiled for linux and macOS.

**Architecture:** A new `cmd/luchta-tsc-worker` entry point drives a read loop over stdin. Each `Run` message is handled in its own goroutine that ports the project's existing `tsc` worker logic (search `tsconfig.build.json` then `tsconfig.json`, clean stale `.d.ts`, create a `Program`, collect diagnostics, emit, report `inputs`/`outputs`) using `tsgo`'s `internal/compiler` + `internal/tsoptions`, with a Yarn-PnP-aware compiler host landed from PR #1966. Protocol I/O, the compile core, and platform utilities live in a new `internal/luchta` package so the `cmd` shell stays thin.

**Tech Stack:** Go 1.26; `github.com/microsoft/typescript-go` internals (`internal/compiler`, `internal/tsoptions`, `internal/execute/tsc`, `internal/vfs/osvfs`, `internal/bundled`, `internal/pnp`, `internal/vfs/pnpvfs`); `encoding/json`, `bufio`; Hereby for builds.

## Global Constraints

- Module path: `github.com/microsoft/typescript-go` (go.mod:1). All imports use this prefix.
- Go version floor: `go 1.26` (go.mod:3).
- Protocol line cap: `MAX_LINE_LENGTH = 1 << 20` (1 MiB) — the stdin reader buffer must allow lines up to this size.
- Protocol JSON is `type`-tagged and **camelCase** (`exitCode`, `resolveTask`, etc.).
- stdout is reserved for protocol JSONL only — compiler/diagnostic text must never reach the real stdout; route it to `Log`. Free-form errors go to stderr.
- Release build flags: `-trimpath -ldflags=-s -w`, `CGO_ENABLED=0`. Do **not** use the `noembed` tag — the default build embeds `lib.d.ts`, making the binary self-contained.
- **Builds/tests in this environment require `GOMODCACHE=/tmp/gomodcache`** prefixed on every `go build`/`go test` command (the default module cache is read-only).
- Target platforms: `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`.
- Out of scope (do not build): source-file cache across runs, reviewdog rdjson, `command`-override parsing, in-memory incremental Program reuse, Windows/linux-arm builds.

### Pinned upstream signatures (reference — copy exactly)

```go
// internal/tsoptions/tsconfigparsing.go:715
type ParseConfigHost interface {
    FS() vfs.FS
    GetCurrentDirectory() string
}

// internal/tsoptions/tsconfigparsing.go:1792
func GetParsedCommandLineOfConfigFile(
    configFileName string,
    options *core.CompilerOptions,
    optionsRaw *collections.OrderedMap[string, any],
    sys ParseConfigHost,
    extendedConfigCache ExtendedConfigCache,
) (*ParsedCommandLine, []*ast.Diagnostic)

// internal/tsoptions/parsedcommandline.go — accessors
func (p *ParsedCommandLine) FileNames() []string
func (p *ParsedCommandLine) CompilerOptions() *core.CompilerOptions
// fields: p.ParsedConfig.ProjectReferences, p.Errors, p.Raw (any)

// internal/execute/tsc/extendedconfigcache.go:15  (zero value usable)
type ExtendedConfigCache struct { /* ... */ }

// internal/compiler/host.go (POST-PnP-merge signature — note the pnpApi param)
func NewCachedFSCompilerHost(
    currentDirectory string,
    fs vfs.FS,
    defaultLibraryPath string,
    extendedConfigCache tsoptions.ExtendedConfigCache,
    pnpApi *pnp.PnpApi,
    trace func(msg *diagnostics.Message, args ...any),
) CompilerHost

// internal/pnp/pnp.go:9 — built-in PnP manifest discovery (walks up from cwd for .pnp.cjs; nil if none)
func InitPnpApi(fs vfs.FS, cwd string) *pnp.PnpApi
// internal/vfs/pnpvfs/pnpvfs.go:24 — wraps an FS to resolve Yarn virtual (zip) paths
func From(fs vfs.FS) vfs.FS

// internal/compiler/program.go:35 & :269
type ProgramOptions struct {
    Host    CompilerHost
    Config  *tsoptions.ParsedCommandLine
    Tracing *tracing.Tracing
    // ... other fields default-zero
}
func NewProgram(opts ProgramOptions) *Program

// internal/compiler/program.go — diagnostics (ctx-taking unless noted)
func (p *Program) GetSyntacticDiagnostics(ctx context.Context, sf *ast.SourceFile) []*ast.Diagnostic
func (p *Program) GetBindDiagnostics(ctx context.Context, sf *ast.SourceFile) []*ast.Diagnostic
func (p *Program) GetSemanticDiagnostics(ctx context.Context, sf *ast.SourceFile) []*ast.Diagnostic
func (p *Program) GetGlobalDiagnostics(ctx context.Context) []*ast.Diagnostic
func (p *Program) GetDeclarationDiagnostics(ctx context.Context, sf *ast.SourceFile) []*ast.Diagnostic
func (p *Program) GetProgramDiagnostics() []*ast.Diagnostic        // no ctx
func (p *Program) GetConfigFileParsingDiagnostics() []*ast.Diagnostic // no ctx

// internal/compiler/program.go:1607,1628
type WriteFile func(fileName string, text string, data *WriteFileData) error
type EmitOptions struct { TargetSourceFile *ast.SourceFile; EmitOnly EmitOnly; WriteFile WriteFile }
func (p *Program) Emit(ctx context.Context, options EmitOptions) *EmitResult
// EmitResult{ EmitSkipped bool; Diagnostics []*ast.Diagnostic; EmittedFiles []string; ... }

// internal/execute/tsc/compile.go:17
type System interface {
    Writer() io.Writer
    FS() vfs.FS
    DefaultLibraryPath() string
    GetCurrentDirectory() string
    WriteOutputIsTTY() bool
    GetWidthOfTerminal() int
    GetEnvironmentVariable(name string) string
    Now() time.Time
    SinceStart() time.Duration
}
// internal/execute/tsc/diagnostics.go:26,30
type DiagnosticReporter = func(*ast.Diagnostic)
func CreateDiagnosticReporter(sys System, w io.Writer, locale locale.Locale, options *core.CompilerOptions) DiagnosticReporter

// internal/vfs/osvfs/os.go:30
func FS() vfs.FS
// internal/bundled — bundled.WrapFS(osvfs.FS()) and bundled.LibPath()
```

---

## Task 1: Land Yarn PnP support (PR #1966) on the fork

Bring PR #1966 onto this fork as an isolated, test-green commit **before** writing any worker code. This is the largest merge-risk step; keep it standalone so conflicts are contained.

**Files:** touches ~76 files; introduces `internal/pnp/`, `internal/vfs/pnpvfs/`, and edits `internal/module/resolver.go`, `internal/compiler/host.go` (adds a `PnpApi()` accessor on the host), `internal/lsp/server.go`, `internal/core/compileroptions.go`, `internal/modulespecifiers/specifiers.go`.

- [ ] **Step 1: Add upstream remote and fetch the PR branch**

```bash
cd /mnt/projects/dobesv/typescript-go
git remote add upstream https://github.com/microsoft/typescript-go.git 2>/dev/null || true
git fetch upstream pull/1966/head:pnp-1966
```
Expected: `git branch` lists `pnp-1966`.

- [ ] **Step 2: Branch and merge the PR onto our work branch**

```bash
git checkout luchta-tsc-worker
git merge --no-ff pnp-1966 -m "Merge Yarn PnP support (PR #1966)"
```
If conflicts: resolve them favoring our fork's surrounding code while keeping the PR's PnP additions intact, then `git add -A && git commit`. Conflicts most likely in `internal/module/resolver.go` and `internal/compiler/host.go`.

- [ ] **Step 3: Build to confirm it compiles**

Run: `go build ./...`
Expected: exits 0, no errors.

- [ ] **Step 4: Run the module-resolution and PnP test suites**

Run: `go test ./internal/module/... ./internal/pnp/... ./internal/vfs/pnpvfs/...`
Expected: PASS. (If the PR shipped compiler baseline tests, also run `go test ./internal/...` and accept the PR's own new baselines.)

- [ ] **Step 5: Record the host PnP integration points for later**

Read and note the exact API the PR added, which Task 7 depends on:
- the `PnpApi()` accessor signature on `compiler.CompilerHost` (in `internal/compiler/host.go`),
- how `internal/lsp/server.go` constructs the PnP API per session and which manifest path it loads,
- the `pnpvfs` wrapper constructor in `internal/vfs/pnpvfs/`.

Run: `grep -rn "PnpApi\|pnpvfs\.\|pnp\.New"  internal/lsp internal/compiler internal/vfs/pnpvfs`
Expected: prints the constructor + accessor lines. Paste them into Task 7's notes before implementing it.

- [ ] **Step 6: Commit (if the merge resolution produced changes beyond the merge commit)**

```bash
git add -A && git commit -m "Resolve PnP merge; tests green" || echo "nothing to commit"
```

---

## Task 2: Protocol types and codec

**Files:**
- Create: `internal/luchta/protocol.go`
- Test: `internal/luchta/protocol_test.go`

**Interfaces:**
- Produces: `type Run struct{ ID, Command, Cwd string; Env map[string]string }`; `type ResolveTask struct{ ID, Name, Command, Package, Cwd string; Scripts []string; Mode string }`; `func DecodeMessage(line []byte) (any, error)` returning `*Run`, `*ResolveTask`, or error; `type Writer struct{}` with `func NewWriter(w io.Writer) *Writer`, `func (*Writer) Log(id, stream, line string)`, `func (*Writer) Done(id string, exitCode int, inputs, outputs []string)`, `func (*Writer) Resolved(id, decision string)`.

- [ ] **Step 1: Write the failing test**

```go
package luchta

import (
	"bytes"
	"strings"
	"testing"
)

func TestDecodeRun(t *testing.T) {
	msg, err := DecodeMessage([]byte(`{"type":"run","id":"pkg#task","command":"","cwd":"packages/pkg","env":{"A":"b"}}`))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	run, ok := msg.(*Run)
	if !ok {
		t.Fatalf("want *Run, got %T", msg)
	}
	if run.ID != "pkg#task" || run.Cwd != "packages/pkg" || run.Env["A"] != "b" {
		t.Fatalf("bad run: %+v", run)
	}
}

func TestDecodeResolveTask(t *testing.T) {
	msg, err := DecodeMessage([]byte(`{"type":"resolveTask","id":"j","name":"build","command":"","package":"@r/a","cwd":"packages/a","scripts":["build"],"mode":"run"}`))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, ok := msg.(*ResolveTask); !ok {
		t.Fatalf("want *ResolveTask, got %T", msg)
	}
}

func TestWriterEmitsTaggedCamelCase(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.Log("id1", "stdout", "hello")
	w.Done("id1", 0, []string{"src/**"}, []string{"dist/a.js"})
	w.Resolved("id2", "accept")
	out := buf.String()
	for _, want := range []string{
		`{"type":"log","id":"id1","stream":"stdout","line":"hello"}`,
		`"type":"done"`, `"exitCode":0`, `"inputs":["src/**"]`, `"outputs":["dist/a.js"]`,
		`{"type":"resolved","id":"id2","result":{"decision":"accept"}}`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
	// each message on its own line
	if lines := strings.Count(strings.TrimSpace(out), "\n"); lines != 2 {
		t.Fatalf("want 3 lines (2 newlines), got %d:\n%s", lines, out)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/luchta/ -run 'Decode|Writer' -v`
Expected: FAIL — `undefined: DecodeMessage` / `NewWriter`.

- [ ] **Step 3: Write minimal implementation**

```go
package luchta

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

type Run struct {
	ID      string            `json:"id"`
	Command string            `json:"command"`
	Cwd     string            `json:"cwd"`
	Env     map[string]string `json:"env"`
}

type ResolveTask struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Command string   `json:"command"`
	Package string   `json:"package"`
	Cwd     string   `json:"cwd"`
	Scripts []string `json:"scripts"`
	Mode    string   `json:"mode"`
}

// DecodeMessage parses one JSONL line into a *Run or *ResolveTask.
func DecodeMessage(line []byte) (any, error) {
	var probe struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(line, &probe); err != nil {
		return nil, fmt.Errorf("malformed message: %w", err)
	}
	switch probe.Type {
	case "run":
		var r Run
		if err := json.Unmarshal(line, &r); err != nil {
			return nil, fmt.Errorf("malformed run: %w", err)
		}
		return &r, nil
	case "resolveTask":
		var r ResolveTask
		if err := json.Unmarshal(line, &r); err != nil {
			return nil, fmt.Errorf("malformed resolveTask: %w", err)
		}
		return &r, nil
	default:
		return nil, fmt.Errorf("unknown message type %q", probe.Type)
	}
}

// Writer serializes protocol responses onto an io.Writer (one JSON object per line).
type Writer struct {
	mu  sync.Mutex
	enc *json.Encoder
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{enc: json.NewEncoder(w)}
}

func (w *Writer) emit(v any) {
	w.mu.Lock()
	defer w.mu.Unlock()
	_ = w.enc.Encode(v) // json.Encoder.Encode appends '\n'
}

func (w *Writer) Log(id, stream, line string) {
	w.emit(map[string]any{"type": "log", "id": id, "stream": stream, "line": line})
}

func (w *Writer) Done(id string, exitCode int, inputs, outputs []string) {
	m := map[string]any{"type": "done", "id": id, "exitCode": exitCode}
	if inputs != nil {
		m["inputs"] = inputs
	}
	if outputs != nil {
		m["outputs"] = outputs
	}
	w.emit(m)
}

func (w *Writer) Resolved(id, decision string) {
	w.emit(map[string]any{"type": "resolved", "id": id, "result": map[string]any{"decision": decision}})
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/luchta/ -run 'Decode|Writer' -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/luchta/protocol.go internal/luchta/protocol_test.go
git commit -m "Add luchta worker protocol codec"
```

---

## Task 3: Output utilities (relativize + clean stale outputs)

Ports `relativizeOutputs` and `cleanOutputs` from the JS worker. `cleanOutputs` deletes `outDir/**/*.d.ts{,.map}` files whose corresponding source under `rootDir` no longer exists.

**Files:**
- Create: `internal/luchta/outputs.go`
- Test: `internal/luchta/outputs_test.go`

**Interfaces:**
- Produces: `func RelativizeOutputs(cwd string, outputs []string) []string`; `func CleanOutputs(cwd, outDir, rootDir string, noEmit bool) error`.

- [ ] **Step 1: Write the failing test**

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/luchta/ -run 'Relativize|CleanOutputs' -v`
Expected: FAIL — `undefined: RelativizeOutputs` / `CleanOutputs`.

- [ ] **Step 3: Write minimal implementation**

```go
package luchta

import (
	"os"
	"path/filepath"
	"strings"
)

// RelativizeOutputs makes absolute output paths relative to cwd (forward slashes).
func RelativizeOutputs(cwd string, outputs []string) []string {
	out := make([]string, 0, len(outputs))
	for _, o := range outputs {
		rel, err := filepath.Rel(cwd, o)
		if err != nil {
			rel = o
		}
		out = append(out, filepath.ToSlash(rel))
	}
	return out
}

// CleanOutputs removes *.d.ts and *.d.ts.map files under outDir whose originating
// source file under rootDir no longer exists. No-op when noEmit is true or outDir
// is absent.
func CleanOutputs(cwd, outDir, rootDir string, noEmit bool) error {
	if noEmit {
		return nil
	}
	absOut := filepath.Join(cwd, outDir)
	absRoot := filepath.Join(cwd, rootDir)
	if _, err := os.Stat(absOut); os.IsNotExist(err) {
		return nil
	}
	return filepath.WalkDir(absOut, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		base := d.Name()
		var stem string
		switch {
		case strings.HasSuffix(base, ".d.ts.map"):
			stem = strings.TrimSuffix(base, ".d.ts.map")
		case strings.HasSuffix(base, ".d.ts"):
			stem = strings.TrimSuffix(base, ".d.ts")
		default:
			return nil
		}
		rel, err := filepath.Rel(absOut, filepath.Join(filepath.Dir(path), stem))
		if err != nil {
			return nil
		}
		// source could be .ts or .tsx
		if fileExists(filepath.Join(absRoot, rel+".ts")) || fileExists(filepath.Join(absRoot, rel+".tsx")) {
			return nil
		}
		return os.Remove(path)
	})
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/luchta/ -run 'Relativize|CleanOutputs' -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/luchta/outputs.go internal/luchta/outputs_test.go
git commit -m "Add output relativize + stale-dts cleanup utilities"
```

---

## Task 4: Compile core (non-PnP)

The heart of the worker: given a package `cwd`, compile `tsconfig.build.json` then `tsconfig.json`, capturing diagnostics, inputs, and outputs. PnP is wired in Task 7; this task uses the OS filesystem only so it can be tested without a PnP manifest.

**Files:**
- Create: `internal/luchta/system.go` (per-run System + ParseConfigHost)
- Create: `internal/luchta/compile.go`
- Test: `internal/luchta/compile_test.go`

**Interfaces:**
- Consumes: `RelativizeOutputs`, `CleanOutputs` (Task 3); the PnP API landed in Task 1 (`pnp.InitPnpApi`, `pnpvfs.From`, the `pnpApi` host-constructor param).
- Produces: `type CompileResult struct{ ExitCode int; Inputs, Outputs []string; Diagnostics string }`; `func CompilePackage(ctx context.Context, cwd string) CompileResult`. A `newRunSystem(cwd string, fsys vfs.FS, libraryPath string, w io.Writer) *runSystem` implementing `tsc.System` and `tsoptions.ParseConfigHost`.

**PnP wiring (folded in from former Task 6 — Task 1's merge made this trivial):** Per compile, build the FS and PnP API directly. `pnp.InitPnpApi` does its own upward search for `.pnp.cjs`, so no custom manifest discovery is needed:
```go
fsys := bundled.WrapFS(osvfs.FS())
pnpApi := pnp.InitPnpApi(fsys, cwd) // nil when no .pnp.cjs above cwd
if pnpApi != nil {
    fsys = pnpvfs.From(fsys)
}
host := compiler.NewCachedFSCompilerHost(cwd, fsys, bundled.LibPath(), extendedConfigCache, pnpApi, nil)
```
Non-PnP packages get `pnpApi == nil` and the plain FS, so the Task 4 tests (no `.pnp.cjs`) still pass. The dedicated PnP resolution test lives in Task 6.

- [ ] **Step 1: Write the failing test**

```go
package luchta

import (
	"context"
	"os"
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
	_ = os.Stat // keep import if unused elsewhere
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/luchta/ -run 'CompilePackage' -v`
Expected: FAIL — `undefined: CompilePackage`.

- [ ] **Step 3: Write the per-run System (`system.go`)**

```go
package luchta

import (
	"io"
	"time"

	"github.com/microsoft/typescript-go/internal/vfs"
)

// runSystem satisfies internal/execute/tsc.System and internal/tsoptions.ParseConfigHost
// for a single package compilation rooted at cwd. Writer routes diagnostic text to Log.
type runSystem struct {
	cwd   string
	fs    vfs.FS
	libs  string
	out   io.Writer
	start time.Time
}

func newRunSystem(cwd string, fsys vfs.FS, libraryPath string, w io.Writer) *runSystem {
	return &runSystem{cwd: cwd, fs: fsys, libs: libraryPath, out: w, start: time.Now()}
}

func (s *runSystem) Writer() io.Writer                       { return s.out }
func (s *runSystem) FS() vfs.FS                              { return s.fs }
func (s *runSystem) DefaultLibraryPath() string              { return s.libs }
func (s *runSystem) GetCurrentDirectory() string             { return s.cwd }
func (s *runSystem) WriteOutputIsTTY() bool                  { return false }
func (s *runSystem) GetWidthOfTerminal() int                 { return 0 }
func (s *runSystem) GetEnvironmentVariable(name string) string { return "" }
func (s *runSystem) Now() time.Time                          { return time.Now() }
func (s *runSystem) SinceStart() time.Duration               { return time.Since(s.start) }
```

> Note: `time.Now()` is used only for diagnostics/timing here, never for protocol identity, so it is safe.

- [ ] **Step 4: Write the compile core (`compile.go`)**

```go
package luchta

import (
	"bytes"
	"context"
	"path/filepath"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/pnp"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
	"github.com/microsoft/typescript-go/internal/vfs/pnpvfs"
)

type CompileResult struct {
	ExitCode    int
	Inputs      []string
	Outputs     []string
	Diagnostics string
}

var tsconfigCandidates = []string{"tsconfig.build.json", "tsconfig.json"}

// CompilePackage replicates the project's tsc worker: for each candidate tsconfig
// that exists, parse it, clean stale outputs, build a Program, collect diagnostics,
// and emit. Returns aggregated inputs/outputs and a non-zero ExitCode on any error.
func CompilePackage(ctx context.Context, cwd string) CompileResult {
	fsys := bundled.WrapFS(osvfs.FS())
	pnpApi := pnp.InitPnpApi(fsys, cwd) // nil when there is no .pnp.cjs above cwd
	if pnpApi != nil {
		fsys = pnpvfs.From(fsys)
	}
	var diagBuf bytes.Buffer
	sys := newRunSystem(cwd, fsys, bundled.LibPath(), &diagBuf)
	extendedConfigCache := &tsc.ExtendedConfigCache{}

	inputs := collections.NewSetWithSizeHint[string](4)
	outputs := []string{}
	exitCode := 0

	for _, name := range tsconfigCandidates {
		configPath := filepath.Join(cwd, name)
		if !fsys.FileExists(configPath) {
			continue
		}
		inputs.Add(name)

		parsed, errs := tsoptions.GetParsedCommandLineOfConfigFile(
			configPath, &core.CompilerOptions{}, nil, sys, extendedConfigCache,
		)
		if len(errs) > 0 {
			reportDiagnostics(sys, &diagBuf, parsed, errs)
			exitCode = 1
			break
		}
		for _, p := range includePatterns(parsed) {
			inputs.Add(p)
		}

		opts := parsed.CompilerOptions()
		if err := CleanOutputs(cwd, outDirOf(opts), rootDirOf(opts), opts.NoEmit.IsTrue()); err != nil {
			diagBuf.WriteString("clean outputs failed: " + err.Error() + "\n")
			exitCode = 1
			break
		}

		host := compiler.NewCachedFSCompilerHost(cwd, fsys, sys.DefaultLibraryPath(), extendedConfigCache, pnpApi, nil)
		program := compiler.NewProgram(compiler.ProgramOptions{Config: parsed, Host: host})

		diags := collectAllDiagnostics(ctx, program)

		emitResult := program.Emit(ctx, compiler.EmitOptions{
			WriteFile: func(fileName, text string, data *compiler.WriteFileData) error {
				outputs = append(outputs, fileName)
				if err := osvfs.FS().WriteFile(fileName, text, false); err != nil {
					return err
				}
				return nil
			},
		})
		diags = append(diags, emitResult.Diagnostics...)

		if len(diags) > 0 {
			reportDiagnostics(sys, &diagBuf, parsed, diags)
			exitCode = 1
			break
		}
	}

	if inputs.Len() == 0 {
		inputs.Add("src/**")
	}
	return CompileResult{
		ExitCode:    exitCode,
		Inputs:      sortedKeys(inputs),
		Outputs:     RelativizeOutputs(cwd, outputs),
		Diagnostics: diagBuf.String(),
	}
}

func collectAllDiagnostics(ctx context.Context, p *compiler.Program) []*ast.Diagnostic {
	var d []*ast.Diagnostic
	d = append(d, p.GetConfigFileParsingDiagnostics()...)
	d = append(d, p.GetSyntacticDiagnostics(ctx, nil)...)
	d = append(d, p.GetProgramDiagnostics()...)
	d = append(d, p.GetBindDiagnostics(ctx, nil)...)
	d = append(d, p.GetGlobalDiagnostics(ctx)...)
	d = append(d, p.GetSemanticDiagnostics(ctx, nil)...)
	d = append(d, p.GetDeclarationDiagnostics(ctx, nil)...)
	return d
}

func reportDiagnostics(sys tsc.System, w *bytes.Buffer, parsed *tsoptions.ParsedCommandLine, diags []*ast.Diagnostic) {
	var opts *core.CompilerOptions
	if parsed != nil {
		opts = parsed.CompilerOptions()
	}
	report := tsc.CreateDiagnosticReporter(sys, w, "", opts) // "" => default (English) locale
	for _, d := range diags {
		report(d)
	}
}

func outDirOf(o *core.CompilerOptions) string {
	if o != nil && o.OutDir != "" {
		return o.OutDir
	}
	return "dist/types"
}

func rootDirOf(o *core.CompilerOptions) string {
	if o != nil && o.RootDir != "" {
		return o.RootDir
	}
	return "src"
}

// includePatterns extracts the tsconfig "include" globs from the raw config,
// defaulting to ["src/**"] (mirrors the JS worker's input reporting).
func includePatterns(parsed *tsoptions.ParsedCommandLine) []string {
	def := []string{"src/**"}
	if parsed == nil || parsed.Raw == nil {
		return def
	}
	var inc any
	switch raw := parsed.Raw.(type) {
	case map[string]any:
		inc = raw["include"]
	case *collections.OrderedMap[string, any]:
		if v, ok := raw.Get("include"); ok {
			inc = v
		}
	}
	arr, ok := inc.([]any)
	if !ok || len(arr) == 0 {
		return def
	}
	out := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		return def
	}
	return out
}

func sortedKeys(s *collections.Set[string]) []string {
	keys := make([]string, 0, s.Len())
	for k := range s.Keys() {
		keys = append(keys, k)
	}
	slicesSortStrings(keys)
	return keys
}
```

> Implementation notes for the engineer:
> - Verify the exact names of the helpers used here against the landed code and adjust if a signature differs: `core.CompilerOptions.OutDir`/`RootDir`/`NoEmit` (NoEmit is a `core.Tristate`; use its `.IsTrue()` method — confirm the method name), `collections.NewSetWithSizeHint`/`Set.Add`/`Set.Len`/`Set.Keys`, and `vfs.FS.WriteFile(path, contents string, writeByteOrderMark bool) error`. If `collections` lacks a `Set`, use a `map[string]struct{}` and sort keys with `sort.Strings`; replace `slicesSortStrings` accordingly (`sort.Strings(keys)`).
> - `parsed.Raw`'s concrete type is either `map[string]any` or `*collections.OrderedMap[string,any]`; the type switch handles both.
> - `WriteFile` writes via `osvfs.FS()` directly (the package's real files), not through the bundled/PnP wrapper, since outputs go to the real disk.

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./internal/luchta/ -run 'CompilePackage' -v`
Expected: PASS. Fix any signature mismatches per the implementation notes until green.

- [ ] **Step 6: Commit**

```bash
git add internal/luchta/system.go internal/luchta/compile.go internal/luchta/compile_test.go
git commit -m "Add tsgo-backed compile core for the worker"
```

---

## Task 5: Worker run loop and binary entry point

Wire the protocol codec to the compile core: read stdin line-by-line, dispatch each message to a goroutine, run the compile, stream diagnostics to `Log`, finish with `Done`. Answer `ResolveTask` with `accept`. Recover panics into `Done{exitCode:1}`.

**Files:**
- Create: `internal/luchta/worker.go`
- Create: `cmd/luchta-tsc-worker/main.go`
- Test: `internal/luchta/worker_test.go`

**Interfaces:**
- Consumes: `DecodeMessage`, `NewWriter`, `*Writer`, `CompilePackage` (Tasks 2, 4).
- Produces: `func Serve(ctx context.Context, in io.Reader, out io.Writer, errw io.Writer) error`.

- [ ] **Step 1: Write the failing test**

```go
package luchta

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"
)

func TestServeRunEmitsLogAndDone(t *testing.T) {
	cwd := t.TempDir()
	writeTsPackage(t, cwd, `{"compilerOptions":{"outDir":"dist","rootDir":"src","module":"nodenext","moduleResolution":"nodenext"}}`,
		"index.ts", "export const x = 1;\n")

	// JSON-encode cwd safely into the Run message.
	in := strings.NewReader(`{"type":"run","id":"t1","command":"","cwd":` + jsonString(cwd) + `,"env":{}}` + "\n")
	var out, errw bytes.Buffer
	if err := Serve(context.Background(), in, &out, &errw); err != nil {
		t.Fatalf("Serve: %v", err)
	}
	s := out.String()
	if !strings.Contains(s, `"type":"done"`) || !strings.Contains(s, `"id":"t1"`) || !strings.Contains(s, `"exitCode":0`) {
		t.Fatalf("missing done: %s", s)
	}
	if !fileExists(filepath.Join(cwd, "dist", "index.js")) {
		t.Fatalf("expected emit")
	}
}

func TestServeResolveTaskAccepts(t *testing.T) {
	in := strings.NewReader(`{"type":"resolveTask","id":"r1","name":"build","command":"","package":"p","cwd":"x","scripts":[],"mode":"run"}` + "\n")
	var out, errw bytes.Buffer
	if err := Serve(context.Background(), in, &out, &errw); err != nil {
		t.Fatalf("Serve: %v", err)
	}
	if !strings.Contains(out.String(), `"decision":"accept"`) {
		t.Fatalf("expected accept: %s", out.String())
	}
}

func TestServeMalformedLineGoesToStderr(t *testing.T) {
	in := strings.NewReader("not json\n")
	var out, errw bytes.Buffer
	if err := Serve(context.Background(), in, &out, &errw); err != nil {
		t.Fatalf("Serve: %v", err)
	}
	if out.Len() != 0 {
		t.Fatalf("malformed input must not write to protocol stdout: %s", out.String())
	}
	if errw.Len() == 0 {
		t.Fatalf("expected stderr diagnostic")
	}
}
```

- [ ] **Step 2: Add the `jsonString` test helper to `protocol_test.go`**

```go
func jsonString(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
```
(Add `"encoding/json"` to that test file's imports.)

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./internal/luchta/ -run 'Serve' -v`
Expected: FAIL — `undefined: Serve`.

- [ ] **Step 4: Write the run loop (`worker.go`)**

```go
package luchta

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
)

const maxLineLength = 1 << 20 // matches luchta MAX_LINE_LENGTH

// Serve reads JSONL messages from in, dispatches each Run/ResolveTask, and writes
// protocol responses to out. Free-form parse errors go to errw (stderr). It returns
// when in reaches EOF, after all in-flight Runs complete.
func Serve(ctx context.Context, in io.Reader, out io.Writer, errw io.Writer) error {
	w := NewWriter(out)
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 0, 64*1024), maxLineLength)

	var wg sync.WaitGroup
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		msg, err := DecodeMessage([]byte(line))
		if err != nil {
			fmt.Fprintf(errw, "luchta-tsc-worker: %v\n", err)
			continue
		}
		switch m := msg.(type) {
		case *ResolveTask:
			w.Resolved(m.ID, "accept")
		case *Run:
			wg.Add(1)
			go func(run *Run) {
				defer wg.Done()
				handleRun(ctx, w, run)
			}(m)
		}
	}
	wg.Wait()
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func handleRun(ctx context.Context, w *Writer, run *Run) {
	defer func() {
		if r := recover(); r != nil {
			w.Log(run.ID, "stderr", fmt.Sprintf("panic: %v", r))
			w.Done(run.ID, 1, nil, nil)
		}
	}()
	res := CompilePackage(ctx, run.Cwd)
	for _, line := range strings.Split(strings.TrimRight(res.Diagnostics, "\n"), "\n") {
		if line != "" {
			w.Log(run.ID, "stdout", line)
		}
	}
	w.Done(run.ID, res.ExitCode, res.Inputs, res.Outputs)
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./internal/luchta/ -run 'Serve' -v`
Expected: PASS.

- [ ] **Step 6: Write the binary entry point (`cmd/luchta-tsc-worker/main.go`)**

```go
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/microsoft/typescript-go/internal/luchta"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := luchta.Serve(ctx, os.Stdin, os.Stdout, os.Stderr); err != nil {
		os.Exit(1)
	}
}
```

- [ ] **Step 7: Build the binary and smoke-test it**

```bash
go build -o ./built/local/luchta-tsc-worker ./cmd/luchta-tsc-worker
printf '{"type":"resolveTask","id":"r","name":"build","command":"","package":"p","cwd":"x","scripts":[],"mode":"run"}\n' | ./built/local/luchta-tsc-worker
```
Expected: prints `{"type":"resolved","id":"r","result":{"decision":"accept"}}`.

- [ ] **Step 8: Commit**

```bash
git add internal/luchta/worker.go internal/luchta/worker_test.go cmd/luchta-tsc-worker/main.go
git commit -m "Add luchta-tsc-worker run loop and binary entry point"
```

---

## Task 6: PnP resolution test (wiring landed in Task 4)

The PnP wiring (`pnp.InitPnpApi` + `pnpvfs.From` + the `pnpApi` host param) was folded into Task 4's `CompilePackage`. This task proves that seam works: when a `.pnp.cjs` exists above the package, the PnP code path is taken, and (where reproducible) a PnP-only dependency resolves. Deep resolution correctness is already covered by PR #1966's own compiler tests (`testdata/tests/cases/compiler/pnp*.ts`); this task guards our worker's use of that machinery.

**Files:**
- Modify: `internal/luchta/compile.go` (extract the FS/PnP construction into a small testable helper)
- Test: `internal/luchta/pnp_test.go`

**Interfaces:**
- Consumes: `pnp.InitPnpApi(fs vfs.FS, cwd string) *pnp.PnpApi`, `pnpvfs.From(fs vfs.FS) vfs.FS`, `bundled.WrapFS`, `osvfs.FS` (all landed in Task 1).
- Produces: `func compilerFS(cwd string) (vfs.FS, *pnp.PnpApi)` — returns the (possibly PnP-wrapped) FS and the PnP API (nil when no `.pnp.cjs` is found above `cwd`).

- [ ] **Step 1: Extract the helper in `compile.go`**

Replace the inline FS/PnP block at the top of `CompilePackage` with a call to a named helper, so the seam is unit-testable:

```go
// compilerFS builds the vfs for compiling under cwd, enabling Yarn PnP when a
// .pnp.cjs manifest exists at or above cwd. Returns the FS and the PnP API (nil
// when not a PnP workspace).
func compilerFS(cwd string) (vfs.FS, *pnp.PnpApi) {
	fsys := bundled.WrapFS(osvfs.FS())
	pnpApi := pnp.InitPnpApi(fsys, cwd)
	if pnpApi != nil {
		fsys = pnpvfs.From(fsys)
	}
	return fsys, pnpApi
}
```

And in `CompilePackage`:

```go
	fsys, pnpApi := compilerFS(cwd)
```

(Add the `vfs` import: `"github.com/microsoft/typescript-go/internal/vfs"`.)

- [ ] **Step 2: Write the failing test**

```go
package luchta

import (
	"os"
	"path/filepath"
	"testing"
)

// Locate the simplest PnP fixture shipped by PR #1966.
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
```

> If `pnp.InitPnpApi` requires the manifest's referenced directories to exist before it returns non-nil (i.e. it validates the manifest, not just its presence), `TestCompilerFSDetectsPnp` will fail. In that case, switch the assertion to construct the manifest scenario the way the PR's `internal/pnp` tests do — find them with `grep -rln "InitPnpApi\|\.pnp\.cjs" internal/pnp` and mirror their setup. Do not weaken the test to assert nothing; if a full setup is infeasible in a unit test, mark this case `t.Skip` with a comment pointing at the PR's compiler-level `pnp*.ts` tests that cover resolution, and keep `TestCompilerFSNoPnp` as the guaranteed gate.

- [ ] **Step 3: Run test to verify it fails**

Run: `GOMODCACHE=/tmp/gomodcache go test ./internal/luchta/ -run 'CompilerFS' -v`
Expected: FAIL — `undefined: compilerFS` (before Step 1 is applied) — then, after Step 1, the detection assertions drive any needed fixup.

- [ ] **Step 4: Make it pass**

Apply Step 1's helper, adjust the detection test per the note if the fixture needs more setup, and re-run.

Run: `GOMODCACHE=/tmp/gomodcache go test ./internal/luchta/ -run 'CompilerFS' -v`
Expected: PASS.

- [ ] **Step 5: Full-package regression**

Run: `GOMODCACHE=/tmp/gomodcache go test ./internal/luchta/...`
Expected: PASS — all earlier tests still green (non-PnP packages get `pnpApi == nil`).

- [ ] **Step 6: Commit**

```bash
git add internal/luchta/compile.go internal/luchta/pnp_test.go
git commit -m "Test Yarn PnP detection in the compile FS seam"
```

---

## Task 7: Cross-compile build task

Add a Hereby task that produces release binaries for all four target platforms.

**Files:**
- Modify: `Herebyfile.mjs` (add a `worker:build` task near the existing `tsgo:build`)

- [ ] **Step 1: Add the build task to `Herebyfile.mjs`**

Insert after the existing `tsgoBuild` task definition:

```javascript
const workerPlatforms = [
    ["linux", "amd64"],
    ["linux", "arm64"],
    ["darwin", "amd64"],
    ["darwin", "arm64"],
];

export const workerBuild = task({
    name: "worker:build",
    description: "Cross-compiles luchta-tsc-worker for linux and macOS.",
    run: async () => {
        for (const [goos, goarch] of workerPlatforms) {
            const out = `./built/worker/${goos}-${goarch}/luchta-tsc-worker`;
            await $({ env: { ...process.env, GOOS: goos, GOARCH: goarch, CGO_ENABLED: "0" } })`go build -trimpath ${["-ldflags=-s -w"]} -o ${out} ./cmd/luchta-tsc-worker`;
        }
    },
});
```

> Match the surrounding file's conventions: reuse the existing `$` import/helper and the `task(...)` factory already used by `tsgoBuild`. If the file builds `ldflags` via a helper (`getReleaseBuildFlags`), prefer reusing it. Do **not** add the `noembed` tag — embedded libs keep each binary self-contained.

- [ ] **Step 2: Run the build task**

Run: `npx hereby worker:build`
Expected: exits 0; creates `built/worker/{linux-amd64,linux-arm64,darwin-amd64,darwin-arm64}/luchta-tsc-worker`.

- [ ] **Step 3: Verify the produced binaries' architectures**

Run: `file built/worker/*/luchta-tsc-worker`
Expected: reports ELF x86-64, ELF aarch64, Mach-O x86_64, Mach-O arm64 respectively.

- [ ] **Step 4: Commit**

```bash
git add Herebyfile.mjs
git commit -m "Add worker:build cross-compile task for linux and macOS"
```

---

## Task 8: End-to-end validation in a real luchta run

Confirm the binary works as a luchta worker against the actual project.

- [ ] **Step 1: Place the binary on PATH and configure luchta**

Copy the host-platform binary somewhere on `PATH` (e.g. `cp built/worker/$(go env GOOS)-$(go env GOARCH)/luchta-tsc-worker ~/.local/bin/`). In the target project's `luchta-config.*`, add a worker: `workers: { tsc: "luchta-tsc-worker" }`, and point the type-check task(s) at that worker.

- [ ] **Step 2: Run a build through luchta**

Run the project's normal luchta build for a couple of packages.
Expected: tasks succeed; diagnostics from any deliberately-broken file appear in luchta's task output; `Done` inputs/outputs feed luchta's cache (a second run is a cache hit when nothing changed).

- [ ] **Step 3: Confirm a failing package reports correctly and the worker survives**

Introduce a type error in one package, run the build.
Expected: that task fails with the diagnostic text; other concurrent tasks still complete; the resident worker process does not crash (no spawn-retry storm in luchta logs).

- [ ] **Step 4: Document usage**

Add a short `cmd/luchta-tsc-worker/README.md` describing the `luchta-config` worker entry, the `tsconfig.build.json`/`tsconfig.json` search behavior, and that modules resolve via Yarn PnP.

```bash
git add cmd/luchta-tsc-worker/README.md
git commit -m "Document luchta-tsc-worker usage"
```

---

## Self-Review

**Spec coverage:**
- PnP support (spec Component 1) → Task 1 + Task 6. ✓
- Protocol layer (Component 2) → Task 2. ✓
- Run handling & concurrency, ResolveTask accept, panic recovery (Component 3) → Task 5. ✓
- Compile core / tsconfig search / cleanOutputs / inputs+outputs (Component 4) → Tasks 3 + 4. ✓
- PnP-aware host (Component 5) → Task 6. ✓
- Cross-compile build (Component 6) → Task 7. ✓
- Testing (unit + integration) → embedded in Tasks 2–6; end-to-end in Task 8. ✓
- Error handling (malformed line → stderr; panic → Done exit 1; missing tsconfig → exit 0 + src/**) → Tasks 4 & 5 tests. ✓

**Known implementation-time confirmations (called out inline, not placeholders):**
- Exact `core.CompilerOptions` field/method names (`OutDir`, `RootDir`, `NoEmit.IsTrue()`), `collections.Set` API, and `vfs.FS.WriteFile` signature — Task 4 Step 4 notes give the fallback (`map[string]struct{}` + `sort.Strings`).
- The PnP wiring shape (FS wrapper vs. host accessor vs. ProgramOptions field) — pinned by Task 1 Step 5 notes, consumed in Task 6.

These are unavoidable: the PnP API does not exist in the tree until Task 1 lands PR #1966, so Task 6 reads the just-landed code rather than guessing its signatures.
