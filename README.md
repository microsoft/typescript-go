# TypeScript 7

> [!NOTE]
> **Parable fork notice.** This is `parable-work/typescript-go`, a fork of
> [microsoft/typescript-go](https://github.com/microsoft/typescript-go) (Project Corsa). It exists
> solely to expose the compiler API to Parable's `psgen` TypeScript reader, which drives the
> compiler in-process from the `parable-platform` Go module. See
> [Fork policy and upgrade procedure](#fork-policy-and-upgrade-procedure) below.

[Not sure what this is? Read the announcement post!](https://devblogs.microsoft.com/typescript/typescript-native-port/)

## Fork policy and upgrade procedure

### What this fork changes

Upstream keeps its entire Go API under `internal/`, which external modules cannot import. This fork
adds a `shim/` re-export layer — and nothing else:

- `shim/<pkg>` packages (`compiler`, `checker`, `ast`, `tsoptions`, `vfs`, `vfs/osvfs`, `bundled`,
  `core`) re-export internal APIs via type aliases (`type Program = compiler.Program`) and var/const
  bindings (`var NewProgram = compiler.NewProgram`). No wrapper logic. The surface is intentionally
  minimal and expands only as `psgen` walker code demands.
- `shim/shim_test.go` is a smoke test that drives the full consumer flow (tsconfig parse → Program
  → type checker → AST walk) through the shim packages only. It is the canary that catches
  upstream API moves during a sync.
- This README section.

**Never touch `internal/`** (or anything else outside the list above) in fork-only commits. That
discipline is what keeps upstream syncs conflict-free.

The module path intentionally remains `github.com/microsoft/typescript-go` so that upstream merges
never conflict on import paths.

### How parable-platform consumes this fork

`parable-platform` imports the shim packages under the upstream module path and redirects the
module to this fork with a `replace` directive pinned to a commit pseudo-version (no tags):

```go
// go.mod
require github.com/microsoft/typescript-go v0.0.0

replace github.com/microsoft/typescript-go => github.com/parable-work/typescript-go v0.0.0-20260610201500-abcdef123456
```

```go
import (
    "github.com/microsoft/typescript-go/shim/compiler"
    "github.com/microsoft/typescript-go/shim/tsoptions"
)
```

To pin a new fork commit: `go mod edit -replace github.com/microsoft/typescript-go=github.com/parable-work/typescript-go@<commit-sha>` then `go mod tidy` (Go resolves the SHA to a pseudo-version).

All Corsa API calls in `parable-platform` stay localized to `utils/psgen/internal/loader/tsreader/`;
nothing else in that repo may import these packages.

### Upgrade procedure (dedicated PR cadence)

Upstream syncs always land as **dedicated PRs** — never mixed with shim surface changes or any
other work:

1. `git remote add upstream https://github.com/microsoft/typescript-go.git` (first time only).
2. `git fetch upstream`.
3. Branch from `main`: `git checkout -b upgrade/upstream-YYYYMMDD`.
4. `git merge upstream/main` (or a specific upstream tag/commit when pinning to a release).
   Conflicts should be rare to nonexistent given the fork-change policy above.
5. Fix any shim compile breaks caused by upstream API moves **in the same PR**, keeping the shim a
   pure re-export layer.
6. Verify: `go build ./...` and `go test ./shim/`.
7. Open the PR, merge to `main`.
8. In `parable-platform`, bump the `replace` pseudo-version to the new fork commit — also as its own
   dedicated PR there — and fix any `tsreader` breaks in that PR.

Shim surface expansions (new aliases/bindings demanded by walker code) follow the same rule in
reverse: they are their own small PRs against fork `main` and never ride along with an upstream
sync.

---

## Preview

A preview build is available on npm as [`@typescript/native-preview`](https://www.npmjs.com/package/@typescript/native-preview).

```sh
npm install @typescript/native-preview
npx tsgo # Use this as you would tsc.
```

A preview VS Code extension is [available on the VS Code marketplace](https://marketplace.visualstudio.com/items?itemName=TypeScriptTeam.native-preview).

To use this, set this in your VS Code settings:

```json
{
    "js/ts.experimental.useTsgo": true
}
```

## What Works So Far?

This is still a work in progress and is not yet at full feature parity with TypeScript. Bugs may exist. Please check this list carefully before logging a new issue or assuming an intentional change.

| Feature | Status | Notes |
|---------|--------|-------|
| Program creation | done | Same files and module resolution as TS 6.0. Not all resolution modes supported yet. |
| Parsing/scanning | done | Exact same syntax errors as TS 6.0 |
| Commandline and `tsconfig.json` parsing | done | Done, though `tsconfig` errors may not be as helpful. |
| Type resolution | done | Same types as TS 6.0. |
| Type checking | done | Same errors, locations, and messages as TS 6.0. Types printback in errors may display differently. |
| JavaScript-specific inference and JSDoc | done | Complete, but intentionally lacking some features. Declaration emit differs greatly, intentionally, to be closer to TS declarations. |
| JSX | done | - |
| Declaration emit | done | - |
| Emit (JS output) | done | - |
| Watch mode | prototype | Watches files and rebuilds, but no incremental rechecking. Not optimized. |
| Build mode / project references | done | - |
| Incremental build | done | - |
| Language service (LSP) | in progress | Nearly all features implemented. |
| API | not ready | - |

Definitions:

 * **done** aka "believed done": We're not currently aware of any deficits or major work left to do. OK to log bugs
 * **in progress**: currently being worked on; some features may work and some might not. OK to log panics, but nothing else please
 * **prototype**: proof-of-concept only; do not log bugs
 * **not ready**: either haven't even started yet, or far enough from ready that you shouldn't bother messing with it yet

## Other Notes

Long-term, we expect that this repo and its contents will be merged into `microsoft/TypeScript`.
As a result, the repo and issue tracker for typescript-go will eventually be closed, so treat discussions/issues accordingly.

For a list of intentional changes with respect to TypeScript 6.0, see CHANGES.md.

## Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit [Contributor License Agreements](https://cla.opensource.microsoft.com).

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

## Trademarks

This project may contain trademarks or logos for projects, products, or services. Authorized use of Microsoft
trademarks or logos is subject to and must follow
[Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/legal/intellectualproperty/trademarks/usage/general).
Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship.
Any use of third-party trademarks or logos are subject to those third-party's policies.
