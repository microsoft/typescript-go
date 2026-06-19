# Design: `luchta-tsc-worker`

**Date:** 2026-06-18
**Status:** Approved (design); pending implementation plan
**Repo:** `github.com/microsoft/typescript-go` (this fork)

## Goal

Produce statically-linked binaries for **linux/{amd64,arm64}** and **darwin/{amd64,arm64}**
that:

1. Speak [luchta](file:///mnt/projects/dobesv/luchta)'s JSONL-over-stdio worker protocol, acting as a
   stay-resident TypeScript compile worker.
2. Type-check and emit each package using tsgo's compiler internals.
3. Resolve modules via **Yarn PnP** (porting/landing
   [microsoft/typescript-go#1966](https://github.com/microsoft/typescript-go/pull/1966)).

This replaces the project's current in-process TypeScript worker (a JS module written against
luchta's TypeScript runtime API) with an external Go binary speaking the raw protocol.

## Background

### luchta worker protocol (summary)

luchta is a Rust build orchestrator for yarn monorepos. It launches stay-resident worker
processes and talks to them over **newline-delimited JSON (JSONL)** on stdio:

- **stdin** (engine → worker): `Run` and `ResolveTask` messages, `type`-tagged, camelCase.
- **stdout** (worker → engine): `Log`, `Done`, `Resolved` responses, `type`-tagged.
- **stderr**: free-form diagnostics, captured by luchta on crash; not part of the protocol.
- Line length cap: `MAX_LINE_LENGTH = 1 << 20` (1 MiB).
- Workers handle **multiple requests concurrently**; responses are correlated by `id`, not by
  arrival order.
- Configured in `luchta-config.*` via `workers: { tsc: "luchta-tsc-worker" }` (command resolved
  on `PATH`). luchta clears the environment and injects only declared + whitelisted vars.

Message shapes (the subset this worker uses):

```
// engine -> worker
Run        { type:"run",         id, command, cwd?, workspace?, env }
ResolveTask{ type:"resolveTask", id, name, command, package, cwd?, scripts[], mode }

// worker -> engine
Log    { type:"log",      id, stream:"stdout"|"stderr", line }
Done   { type:"done",     id, exitCode, inputs?, outputs? }
Resolved{ type:"resolved", id, result:{ decision:"accept"|"prune"|"reject"|"modify", ... } }
```

`Done.inputs`/`Done.outputs`, when present, replace the task's declared cache input/output
patterns for subsequent runs — this is how the worker feeds luchta's build cache.

### Behavior being ported

The existing JS worker compiles **one package per task** (`cwd` = package dir) with plain
`ts.createProgram` (not `tsc -b`). For each of `tsconfig.build.json` then `tsconfig.json`:

- skip if the file does not exist in `cwd`;
- record the tsconfig filename and its `include` patterns (default `src/**`) as **inputs**;
- **clean stale outputs**: delete `outDir/**/*.d.ts{,.map}` (defaults `outDir=dist/types`,
  `rootDir=src`) that have no matching source file, unless `noEmit`;
- create a `Program`, collect semantic/declaration/syntactic/global diagnostics, and `emit`,
  recording every written file as an **output**;
- on any diagnostic or emit failure → print formatted diagnostics and fail.

It builds **both** tsconfigs when both exist, returns `{exitCode, inputs, outputs}` (outputs
relativized to `cwd`), and shares a `WeakRef` source-file cache across compilations (to avoid
re-parsing `lib.d.ts`). The current worker also emits reviewdog rdjson — **out of scope here**.

## Design

### Component 1 — Yarn PnP integration (prerequisite)

Cherry-pick / rebase **PR #1966** onto this fork as an **isolated first step**, with the test
suite green before any worker code is written. It brings:

- `internal/pnp/` — parses `.pnp.cjs` / `.pnp.data.json` (no yarn binary invoked);
- `internal/vfs/pnpvfs/` — zip-backed virtual path filesystem with fallback to standard VFS;
- hooks in `internal/module/resolver.go` (PnP paths tried before standard `node_modules`);
- per-`Host` `PnpApi()` accessor (the reviewer-mandated design enabling per-project isolation
  without global state).

**Risk:** 76-file PR; may conflict with this fork's commits. Mitigation: land it on its own
branch/commit, run the full suite, and only then build on top. This is the single largest source
of effort and merge risk in the project.

### Component 2 — Protocol layer

Go structs mirroring the Rust protocol types, `type`-tagged and camelCase via `encoding/json`:
incoming `Run` / `ResolveTask`; outgoing `Log` / `Done` / `Resolved`.

- A read loop over stdin using `bufio.Scanner` with a **1 MiB buffer cap** matching
  `MAX_LINE_LENGTH`.
- A **mutex-serialized stdout writer** so concurrent tasks' JSONL lines never interleave.
- Decode dispatches on the `type` discriminator. Malformed lines → write to **stderr** and
  continue (do not crash the resident worker).

### Component 3 — Run handling & concurrency

- Each `Run` is handled in its own goroutine; luchta may issue concurrent Runs to one resident
  process.
- A `recover()` per goroutine converts a panic into `Done{exitCode:1}` plus a diagnostic `Log`,
  so one failing task never kills the worker.
- `ResolveTask` → always respond `Resolved{decision:"accept"}` for the MVP.
- All compiler/diagnostic output is routed through the `Log` writer; the process never writes
  compiler output to the real stdout (which is reserved for the protocol).

### Component 4 — Compile core (port of `compileTsconfig` / `tsc`)

Per Run, with `cwd` = package dir:

1. `inputs = {}`, `outputs = {}`.
2. For `tsconfigFile` in `["tsconfig.build.json", "tsconfig.json"]`:
   - if `cwd/tsconfigFile` does not exist → skip;
   - `inputs += tsconfigFile`;
   - parse the tsconfig (tsgo config parsing); `inputs +=` its `include` patterns, or `src/**`
     if none;
   - **cleanOutputs** (skip if `noEmit`): for each `outDir/**/*.d.ts` and `*.d.ts.map`, map back
     to the expected source under `rootDir` (`.d.ts`/`.d.ts.map` → `.ts`/`.tsx`); delete the
     output if no source exists;
   - create a `Program` (rootNames + options + projectReferences + config diagnostics) with the
     PnP-aware host from Component 5;
   - collect semantic, declaration, syntactic, and global diagnostics;
   - `emit`, recording every written path into `outputs`;
   - on emit failure or any diagnostic → format diagnostics to `Log`, set non-zero `exitCode`,
     stop.
3. if `inputs` is empty → `inputs += src/**`.
4. Return `Done{exitCode, inputs, outputs}` with `outputs` relativized to `cwd`.

Two small utilities ride along: `cleanOutputs` (the `clean-dest` logic) and `relativizeOutputs`.

### Component 5 — PnP-aware compiler host

Per Run, locate the nearest PnP manifest upward from `cwd` (typically the repo root) and build a
tsgo `Host` with `PnpApi()` attached, so the `Program` resolves workspace and external
dependencies through PnP rather than `node_modules`. The parsed manifest is **cached by path**
across Runs (same repo root each time; cheap and safe). Cache access is mutex-guarded for
concurrent Runs.

### Component 6 — Build & release

Add a Hereby task that iterates the platform matrix and runs
`go build ./cmd/luchta-tsc-worker` with `GOOS`/`GOARCH`, `CGO_ENABLED=0`, and release flags
`-trimpath -ldflags=-s -w`, emitting static binaries to `built/worker/<os>-<arch>/`. No npm
packaging — the binary is placed on `PATH` and referenced from `luchta-config`. Matrix:
`linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`.

## Testing

- **Unit**
  - Protocol JSON round-trips against the documented luchta wire examples (`run`, `resolveTask`,
    `log`, `done`, `resolved`).
  - `cleanOutputs`: stale `.d.ts`/`.d.ts.map` removed, live ones kept, `noEmit` respected.
  - `relativizeOutputs`.
- **Integration**
  - A fixture PnP monorepo (`.pnp.data.json` + two packages, one depending on the other); feed a
    `Run` as JSONL on stdin; assert the `Log`/`Done` sequence, expected `inputs`/`outputs`/
    `exitCode`, and emitted files on disk.
  - A deliberately-broken package fixture → non-zero `exitCode` + diagnostic `Log`, and the
    resident worker survives to handle a subsequent good `Run`.
  - Reuse PR #1966's resolution baselines for PnP correctness.

## Error handling

- Malformed protocol line → stderr, skip, continue.
- Panic during a Run → `recover()` → `Done{exitCode:1}` + diagnostic `Log`.
- Neither tsconfig present → `exitCode:0`, `inputs:["src/**"]`, no outputs (matches existing
  no-op-success behavior).

## Sequencing

1. Land PnP (PR #1966) on the fork; full test suite green.
2. Protocol layer: structs, stdin read loop, serialized stdout writer.
3. Compile core against tsgo internals (no PnP yet) → green on a non-PnP fixture.
4. Wire the PnP-aware host (Component 5).
5. Cross-compile Hereby task (Component 6).
6. End-to-end test inside a real luchta run.

## Out of scope (deferrable behind the same protocol)

- Source-file cache across runs (start without; measure; add if slow).
- reviewdog rdjson output.
- `command`-override path (worker currently always does the tsconfig-search behavior; luchta
  sends no `command`).
- Incremental in-memory `Program` reuse.
- Windows and linux/arm builds.

## Open risks

- **PR #1966 rebase conflicts** onto this fork (largest).
- Exact tsgo internal APIs for programmatic config parsing, `Program` creation, and capturing
  emitted output paths — to be pinned during planning.
- Concurrency safety when multiple Runs build simultaneously and share the cached PnP manifest /
  any shared VFS caches.
