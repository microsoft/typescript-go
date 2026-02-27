# Writing, Running, and Debugging Compiler Tests and Fourslash Tests

This guide covers the complete testing workflow for the typescript-go repository, including compiler tests (type-checking, emit, diagnostics) and fourslash tests (language server features like completions, hover, go-to-definition).

---

## 1. Compiler Tests

Compiler tests validate the TypeScript compiler's behavior: diagnostics, JavaScript emit, source maps, type/symbol baselines, and more. Each test is a `.ts` or `.tsx` file that the test runner compiles, then compares output against stored baselines.

### 1.1 Where Test Files Live

| Path | Purpose |
|------|---------|
| `testdata/tests/cases/compiler/` | Regression tests (local to this repo) |
| `testdata/tests/cases/conformance/` | Conformance tests (local to this repo) |
| `_submodules/TypeScript/tests/cases/compiler/` | Submodule tests from upstream TypeScript |
| `_submodules/TypeScript/tests/cases/conformance/` | Submodule conformance tests from upstream |

### 1.2 Writing a New Compiler Test

A compiler test is just a `.ts` or `.tsx` file — no Go code needed. Place it in `testdata/tests/cases/compiler/` for regression tests or `testdata/tests/cases/conformance/<subdir>/` for conformance tests.

#### Simple single-file test

```typescript
// testdata/tests/cases/compiler/myNewTest.ts
const x: number = "hello"; // expect type error
```

#### Using compiler option directives

Set compiler options with `// @option: value` comment directives at the top of the file:

```typescript
// @target: es2020
// @strict: true
// @declaration: true
// @jsx: react
// @noEmit: true
const x: number = 42;
```

#### Multi-file test

Use `// @filename:` directives to define multiple files in one test:

```typescript
// @target: es2015
// @module: commonjs

// @filename: /src/utils.ts
export function greet(name: string): string {
    return `Hello, ${name}`;
}

// @filename: /src/main.ts
import { greet } from "./utils";
const msg: number = greet("world"); // type error
```

#### Generating test variations

Options that affect program behavior (those marked with `AffectsProgramStructure`, `AffectsEmit`, `AffectsModuleResolution`, `AffectsDiagnostics`, etc.) can specify multiple values to generate separate sub-test configurations:

```typescript
// @target: es2015, esnext
// @module: commonjs, esnext
// @strict: true, false
export const x = 1;
```

This generates a sub-test for each combination, with names like `myTest.ts (target=es2015,module=commonjs,strict=true)`.

Use `*` to test all valid values for an option, and `!value` to exclude specific ones:

```typescript
// @module: *
// @target: *, !es3
```

#### Symlink tests

Use `// @link:` to create symlinks in the virtual filesystem:

```typescript
// @link: /src -> /node_modules/mylib
```

#### Other directives

- `// @currentDirectory: /custom/path` — Set the working directory
- `// @noImplicitReferences` — Don't auto-include referenced files
- `// @lib: es2020,dom` — Specify lib files

### 1.3 Running Compiler Tests

#### Via hereby (recommended for full runs)

```bash
npx hereby test                          # Run all tests (compiler + fourslash + others)
npx hereby test --tests "myNewTest"      # Run tests matching a pattern (passed as -run flag)
```

#### Via Go directly (faster for specific tests)

```bash
# Run all local compiler tests
go test ./internal/testrunner/ -run TestLocal

# Run a specific test by name
go test ./internal/testrunner/ -run 'TestLocal/compiler/myNewTest.ts'

# Run submodule tests (from upstream TypeScript)
go test ./internal/testrunner/ -run 'TestSubmodule/compiler/someTest.ts'

# Run with verbose output
go test ./internal/testrunner/ -run 'TestLocal/compiler/myNewTest.ts' -v
```

The test entry points in `internal/testrunner/compiler_runner_test.go` are:
- `TestLocal` — runs tests from `testdata/tests/cases/` (both `compiler/` and `conformance/`)
- `TestSubmodule` — runs tests from `_submodules/TypeScript/tests/cases/` and generates diff baselines

#### What happens during a test run

For each test file, the runner:
1. Parses directives (`// @option:`, `// @filename:`, etc.)
2. Generates configurations for each option variation
3. For each configuration, runs these parallel sub-tests:
   - `error` — Verifies diagnostics against `.errors.txt` baseline
   - `output` — Verifies JavaScript emit against `.js` baseline
   - `sourcemap` — Verifies source map output
   - `sourcemap record` — Verifies source map record
   - `union ordering` — Validates AST union type ordering
   - `source file parent pointers` — Validates AST structure integrity

### 1.4 Baseline System

Baselines are the expected output files that test results are compared against.

| Directory | Purpose |
|-----------|---------|
| `testdata/baselines/reference/` | Golden/expected baselines (committed to repo) |
| `testdata/baselines/local/` | Generated during test runs (not committed) |

#### Baseline file types

| Extension | Content |
|-----------|---------|
| `.errors.txt` | Diagnostic error messages |
| `.js` | Emitted JavaScript |
| `.d.ts` | Declaration output |
| `.symbols` | Symbol information |
| `.types` | Type information |
| `.sourcemap.txt` | Source map output |
| `.trace.json` | Trace output |

#### Accepting baselines

After running tests and reviewing the differences:

```bash
# View baseline diffs (requires DIFF env var to be set to a diff tool)
DIFF=code npx hereby diff

# Accept all new baselines (copies local/ → reference/)
npx hereby baseline-accept
```

The `baseline-accept` task:
1. Copies all files from `local/` to `reference/` (excluding `.delete` files)
2. Deletes reference files that have corresponding `.delete` markers in `local/`

#### Baseline tracking

Every test package that uses baselines must include a `TestMain` with `baseline.Track()`:

```go
func TestMain(m *testing.M) {
    defer baseline.Track()()
    m.Run()
}
```

When running via `npx hereby test`, the system:
1. Sets the `TSGO_BASELINE_TRACKING_DIR` environment variable to a temp directory
2. Each test records which baselines it uses
3. After all tests complete, detects unused baselines
4. Creates `.delete` markers for unused baselines
5. Reports an error if unused baselines are found

The `--dirty` flag skips this tracking and also skips clearing `local/` before tests.

### 1.5 Debugging Compiler Tests

#### Examine baseline diffs

```bash
# Set your diff tool
export DIFF=code  # or vimdiff, meld, etc.
npx hereby diff
```

Or manually inspect:
```bash
# View generated output
cat testdata/baselines/local/compiler/myTest.errors.txt

# Compare with expected
diff testdata/baselines/reference/compiler/myTest.errors.txt \
     testdata/baselines/local/compiler/myTest.errors.txt
```

#### Run with verbose output

```bash
go test ./internal/testrunner/ -run 'TestLocal/compiler/myNewTest.ts' -v
```

#### Debug with Delve

Build tests with debug flags disabled so Delve can step through:

```bash
# Build with debug flags (disables optimizations and inlining)
npx hereby test --debug --tests "myNewTest"
```

This passes `-gcflags=all=-N -l` to `go build`, which:
- `-N` disables optimizations
- `-l` disables inlining

Then attach Delve to step through code.

Alternatively, use Delve directly:
```bash
dlv test ./internal/testrunner/ -- -test.run 'TestLocal/compiler/myNewTest.ts'
```

---

## 2. Fourslash Tests

Fourslash tests validate language server (LSP) features: completions, hover/quick info, go-to-definition, find references, rename, code fixes, formatting, and more. They're Go test files that set up TypeScript source with position markers, then verify LSP responses.

### 2.1 Where Test Files Live

| Path | Purpose |
|------|---------|
| `internal/fourslash/tests/*.go` | Hand-written fourslash tests |
| `internal/fourslash/tests/manual/*.go` | Hand-written tests (manual subdirectory) |
| `internal/fourslash/tests/gen/*.go` | Auto-generated from upstream TypeScript fourslash tests |
| `internal/fourslash/` | Test harness and utilities |
| `internal/fourslash/tests/util/` | Shared test constants (`DefaultCommitCharacters`, etc.) |

**Key difference**: Generated tests in `gen/` use `fourslash.SkipIfFailing(t)` for tests that are known to not yet work. Hand-written tests should always pass.

### 2.2 Writing a New Fourslash Test

Create a Go test file in `internal/fourslash/tests/`. The file uses the `fourslash_test` package.

#### Minimal template

```go
package fourslash_test

import (
    "testing"

    "github.com/microsoft/typescript-go/internal/fourslash"
    "github.com/microsoft/typescript-go/internal/testutil"
)

func TestMyFeature(t *testing.T) {
    t.Parallel()
    defer testutil.RecoverAndFail(t, "Panic on fourslash test")
    const content = `
var x/*marker1*/ = 42;
`
    f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
    defer done()
    f.VerifyQuickInfoAt(t, "marker1", "var x: number", "")
}
```

#### Real-world example: Quick Info

```go
func TestBasicQuickInfo(t *testing.T) {
    t.Parallel()
    defer testutil.RecoverAndFail(t, "Panic on fourslash test")
    const content = `
/**
 * Some var
 */
var someVar/*1*/ = 123;

/**
 * Other var
 * See {@link someVar}
 */
var otherVar/*2*/ = someVar;
`
    f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
    defer done()
    f.VerifyQuickInfoAt(t, "1", "var someVar: number", "Some var")
    f.VerifyQuickInfoAt(t, "2", "var otherVar: number",
        "Other var\nSee [someVar](file:///basicQuickInfo.ts#4,5-4,12)")
}
```

#### Real-world example: Editing and Completions

```go
func TestBasicEdit(t *testing.T) {
    t.Parallel()
    defer testutil.RecoverAndFail(t, "Panic on fourslash test")
    const content = `export {};
interface Point {
    x: number;
    y: number;
}
declare const p: Point;
p/*a*/`
    f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
    defer done()
    f.GoToMarker(t, "a")
    f.Insert(t, ".")
    f.GoToEOF(t)
    f.VerifyCompletions(t, nil, &fourslash.CompletionsExpectedList{
        IsIncomplete: false,
        ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
            CommitCharacters: &DefaultCommitCharacters,
        },
        Items: &fourslash.CompletionsExpectedItems{
            Exact: []fourslash.CompletionsExpectedItem{
                &lsproto.CompletionItem{
                    Label:    "x",
                    Kind:     new(lsproto.CompletionItemKindField),
                    SortText: new(string(ls.SortTextLocationPriority)),
                },
                "y",
            },
        },
    })
}
```

#### Marker syntax

Markers define cursor positions and text ranges in the test content:

| Syntax | Description | Example |
|--------|-------------|---------|
| `/*name*/` | Named position marker | `var x/*pos*/ = 1;` |
| `/*1*/`, `/*2*/` | Numbered markers | `foo(/*1*/, /*2*/)` |
| `[|text|]` | Range marker (selects text) | `[|let x: number|]` |

#### Multi-file tests

Use `// @Filename:` (capital F) to define multiple files:

```go
const content = `
// @Filename: /src/utils.ts
export function greet(name: string) { return name; }

// @Filename: /src/main.ts
import { greet } from "./utils";
greet(/*marker*/"world");
`
```

#### Setting compiler options

Embed a `tsconfig.json` file or use directive comments:

```go
const content = `
// @Filename: /tsconfig.json
{ "compilerOptions": { "strict": true, "target": "es2020" } }

// @Filename: /src/test.ts
const x/*1*/ = 42;
`
```

### 2.3 Verification Methods (Common API)

The `fourslash.FourslashTest` type (variable `f`) provides these verification methods:

#### Quick Info / Hover
```go
f.VerifyQuickInfoAt(t, "marker", "var x: number", "documentation text")
f.VerifyBaselineHover(t)  // generates baseline file
```

#### Completions
```go
f.VerifyCompletions(t, "marker", &fourslash.CompletionsExpectedList{
    IsIncomplete: false,
    ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
        CommitCharacters: &DefaultCommitCharacters,
        EditRange:        Ignored,
    },
    Items: &fourslash.CompletionsExpectedItems{
        Includes: []fourslash.CompletionsExpectedItem{
            &lsproto.CompletionItem{Label: "myVar"},
        },
        // Or use Exact for exact match:
        // Exact: []fourslash.CompletionsExpectedItem{"x", "y"},
    },
})
```

Import the test utilities for shared constants:
```go
import . "github.com/microsoft/typescript-go/internal/fourslash/tests/util"
// Provides: DefaultCommitCharacters, Ignored, CompletionGlobalThisItem, etc.
```

#### Navigation
```go
f.VerifyBaselineGoToDefinition(t)           // baseline-based
f.VerifyBaselineGoToTypeDefinition(t)
f.VerifyBaselineGoToImplementation(t)
```

#### References and Rename
```go
f.VerifyBaselineFindAllReferences(t)
f.VerifyBaselineRename(t)
```

#### Diagnostics
```go
f.VerifyNoErrors(t)
f.VerifyErrorExistsBetweenMarkers(t, "start", "end")
f.VerifyBaselineNonSuggestionDiagnostics(t)
```

#### Signature Help
```go
f.VerifyBaselineSignatureHelp(t)
f.VerifyNoSignatureHelp(t)
```

#### Editing (simulating user actions)
```go
f.GoToMarker(t, "marker")    // move cursor to marker position
f.Insert(t, ".")             // type text at cursor
f.Backspace(t, 3)            // delete 3 characters before cursor
f.DeleteAtCaret(t, 5)        // delete 5 characters after cursor
f.Paste(t, "new text")       // paste text
f.Replace(t, start, len, "replacement")
f.GoToEOF(t)                 // move to end of file
f.GoToFile(t, "/src/main.ts") // switch to another file
```

#### Other LSP Features
```go
f.VerifyBaselineDocumentHighlights(t)
f.VerifyBaselineDocumentSymbol(t)
f.VerifyBaselineCallHierarchy(t)
f.VerifyBaselineInlayHints(t)
f.VerifyBaselineSelectionRanges(t)
f.VerifyBaselineClosingTags(t)
f.FormatDocument(t, "/test.ts")
f.VerifyOrganizeImports(t, expectedContent, actionKind, prefs)
```

### 2.4 Running Fourslash Tests

```bash
# Run all fourslash tests
go test ./internal/fourslash/tests/...

# Run a specific test
go test ./internal/fourslash/tests/... -run TestBasicQuickInfo

# Run with pattern matching
go test ./internal/fourslash/tests/... -run "TestCompletion.*"

# Run with verbose output
go test ./internal/fourslash/tests/... -run TestBasicQuickInfo -v

# Via hereby (runs all tests including fourslash)
npx hereby test --tests "TestBasicQuickInfo"
```

### 2.5 Fourslash Baselines

Fourslash tests that use `VerifyBaseline*` methods generate baselines under:

```
testdata/baselines/reference/fourslash/<command>/
```

Where `<command>` is one of: `quickInfo`, `signatureHelp`, `goToDefinition`, `goToType`, `goToImplementation`, `findAllReferences`, `documentHighlights`, `findRenameLocations`, `callHierarchy`, `Code Lenses`, `Document Symbols`, `Inlay Hints`, etc.

File extensions vary by command:
- `.baseline` — quickInfo, signatureHelp, diagnostics, etc.
- `.baseline.jsonc` — most other features
- `.baseline.md` — auto imports
- `.callHierarchy.txt` — call hierarchy

Accept baselines the same way as compiler tests:
```bash
npx hereby baseline-accept
```

### 2.6 Debugging Fourslash Tests

#### Run with verbose output
```bash
go test ./internal/fourslash/tests/... -run TestMyFeature -v
```

#### Debug with Delve
```bash
dlv test ./internal/fourslash/tests/... -- -test.run TestMyFeature
```

Or use the `--debug` flag with hereby:
```bash
npx hereby test --debug --tests "TestMyFeature"
```

#### Common issues

| Problem | Cause | Fix |
|---------|-------|-----|
| "marker not found" | Misspelled marker name | Check marker names are exact (case-sensitive) |
| Unexpected completions | Missing `EditRange: Ignored` | Add `EditRange: Ignored` to `ItemDefaults` if you don't care about edit ranges |
| Test panics | Missing `defer done()` | Always call `defer done()` after `NewFourslash` |
| Baseline mismatch | Changed LSP output | Review diff and accept with `npx hereby baseline-accept` |

### 2.7 Generated vs. Hand-Written Tests

Generated tests (in `gen/`) are auto-converted from the upstream TypeScript fourslash test suite using the script at `internal/fourslash/_scripts/convertFourslash.mts`. They:
- Use `fourslash.SkipIfFailing(t)` for tests that don't pass yet
- Should not be manually edited (they'll be overwritten on regeneration)
- Provide coverage for ported TypeScript behavior

Hand-written tests (directly in `internal/fourslash/tests/` or in `manual/`):
- Must always pass (no `SkipIfFailing`)
- Test specific behaviors, edge cases, or new features
- Are the right place for custom regression tests

---

## 3. General Testing Infrastructure

### 3.1 Key hereby Commands

| Command | Description |
|---------|-------------|
| `npx hereby test` | Run all tests with baseline tracking |
| `npx hereby test --tests "Pattern"` | Run tests matching `-run=Pattern` |
| `npx hereby test --debug` | Build with `-gcflags=all=-N -l` for Delve debugging |
| `npx hereby test --dirty` | Skip clearing local baselines and baseline tracking |
| `npx hereby test --race` | Enable Go race detector |
| `npx hereby test --coverage` | Generate coverage profiles |
| `npx hereby baseline-accept` | Accept local baselines as new reference |
| `npx hereby diff` | Diff baselines with tool from `$DIFF` env var |
| `npx hereby format` | Format code (uses dprint) |
| `npx hereby lint` | Run linters (uses golangci-lint) |

### 3.2 Environment Variables

| Variable | Purpose |
|----------|---------|
| `TSGO_BASELINE_TRACKING_DIR` | Directory for baseline usage tracking files |
| `TS_TEST_PROGRAM_SINGLE_THREADED` | Set to `false` for parallel test programs |
| `DIFF` | Path to diff tool for `npx hereby diff` |
| `TSGO_HEREBY_RACE` | Enable race detector via env var |
| `TSGO_HEREBY_NOEMBED` | Skip embedding bundled assets |

### 3.3 Typical Workflow

#### Adding a new compiler test
1. Create `testdata/tests/cases/compiler/myTest.ts` with test code and directives
2. Run: `go test ./internal/testrunner/ -run 'TestLocal/compiler/myTest.ts'`
3. Review generated baselines in `testdata/baselines/local/compiler/`
4. Accept: `npx hereby baseline-accept`

#### Adding a new fourslash test
1. Create `internal/fourslash/tests/myTest_test.go` with the test function
2. Run: `go test ./internal/fourslash/tests/... -run TestMyTest`
3. Review any generated baselines in `testdata/baselines/local/fourslash/`
4. Accept: `npx hereby baseline-accept`

#### Investigating a test failure
1. Run the specific test: `go test ./internal/testrunner/ -run 'TestLocal/compiler/failingTest.ts' -v`
2. Check baseline diffs: `diff testdata/baselines/reference/compiler/failingTest.errors.txt testdata/baselines/local/compiler/failingTest.errors.txt`
3. If the new output is correct, accept: `npx hereby baseline-accept`
4. If not, fix the code and re-run

### 3.4 Key Source Files

| File | Purpose |
|------|---------|
| `internal/testrunner/compiler_runner.go` | Compiler test runner |
| `internal/testrunner/compiler_runner_test.go` | `TestLocal` and `TestSubmodule` entry points |
| `internal/testrunner/test_case_parser.go` | Parses `// @` directives from test files |
| `internal/fourslash/fourslash.go` | Fourslash test harness (all `Verify*` methods) |
| `internal/fourslash/test_parser.go` | Parses markers (`/*name*/`, `[|range|]`) |
| `internal/fourslash/baselineutil.go` | Fourslash baseline file generation |
| `internal/fourslash/tests/util/` | Shared test constants (`DefaultCommitCharacters`, etc.) |
| `internal/testutil/baseline/baseline.go` | Core baseline comparison and tracking |
| `Herebyfile.mjs` | Build/test task definitions |
