# Issue #2278 investigation: `@stylexjs/babel-plugin` `Options` import

## Summary

The reported failure was reproducible in a December 2025 `tsgo` snapshot, but is **not reproducible in the current local checkout**.

The root cause was **not** `skipLibCheck` itself. `skipLibCheck` only made the real package usable enough to expose a second bug:

- the package's `.d.ts` file is structurally invalid because it combines `export =` with other exported declarations
- older `tsgo` then also incorrectly reported `TS2305` for `import { Options } ...`
- stable TypeScript 5.9 did **not** report that `TS2305`

The most likely root cause was a mismatch in CommonJS `export =` merging:

- TypeScript merges supplemental exports from the original module symbol into the resolved `export =` target
- older `tsgo` resolved the `export =` target, but then looked up named exports on the wrong symbol, so type-only exports like `Options` were not found

That behavior appears to have already been fixed in `tsgo` by commit `418d206bc` (`Fix named export aliases merging with \`export =\` (#3157)`).

## Exact repros used

### Real package repro

Directory used locally:

- `testdata/repros/issue2278-stylex-options/`

Contents:

- dependency: `@stylexjs/babel-plugin@0.16.2`
- `tsconfig.json`: `{ "compilerOptions": { "skipLibCheck": true } }`
- `stylex.ts`: `import { Options } from "@stylexjs/babel-plugin";`

Observed results:

- `typescript@5.9.3`: **no errors**
- current local `tsgo`: **no errors**
- `tsgo` at commit `06abb77cb` (Dec 2025 snapshot): **TS2305**

### Minimal synthetic repro

The compiler regression test added in this change is:

- `testdata/tests/cases/compiler/stylexOptionsExportEqualsMerging.ts`

It models the same important shape:

- package with `types: "./lib/index.d.ts"`
- declaration file that has both:
  - `export type Options = SharedOptions`
  - `export = exported`
- consuming file imports `{ Options }`

This is enough to reproduce the old `tsgo` bug without needing the published package.

Observed results:

- `typescript@5.9.3` with `skipLibCheck: true`: **no errors**
- current local `tsgo`: **no errors**
- `tsgo` at commit `06abb77cb`: **TS2305**

Without `skipLibCheck`, stable TypeScript reports only the declaration-file error (`TS2309`) for the invalid `export =` shape, while the older `tsgo` reported both `TS2309` and the extra incorrect `TS2305`.

## Published package shape

Published metadata for `@stylexjs/babel-plugin@0.16.2`:

- `main`: `lib/index.js`
- `types`: `./lib/index.d.ts`
- no `"exports"` field

Relevant declaration from published `lib/index.d.ts`:

```ts
import type { StyleXOptions } from './utils/state-manager';
export type Options = StyleXOptions;
...
declare const $$EXPORT_DEFAULT_DECLARATION$$: StyleXTransformObj;
export = $$EXPORT_DEFAULT_DECLARATION$$;
```

So `Options` is exported as a **type export on the original module symbol**, while the module also has an `export =` assignment.

## Relevant `tsgo` codepaths

### Older failing path

In the December 2025 snapshot:

- `internal/checker/checker.go`
  - `getExternalModuleMember`
  - `getExportOfModule`
  - `resolveExternalModuleSymbol`
  - `getExportsOfModuleWorker`

Problematic behavior:

1. `getExternalModuleMember` resolved the module through `resolveESModuleSymbol`
2. for `export =` modules, `resolveExternalModuleSymbol` returned the `export =` target
3. then `getExternalModuleMember` called:

   - `getExportOfModule(targetSymbol, nameText, ...)`

4. but `targetSymbol` did not include the supplemental type-only export `Options`
5. result: `symbolFromModule == nil`, then `errorNoModuleMemberSymbol(...)` reported `TS2305`

Specific older lines in the Dec 2025 checkout:

- `internal/checker/checker.go:14282`
- `internal/checker/checker.go:15070-15077`
- `internal/checker/checker.go:15717-15723`

### Current local checkout

The current checkout contains the fix in:

- `internal/checker/checker.go:14228-14233`
- `internal/checker/checker.go:15717-15727`

Notable change:

- `getExternalModuleMember` now switches from `targetSymbol` back to `moduleSymbol` when checking named exports for `export =` modules
- `getExportsOfModuleWorker` also explicitly preserves supplemental type/namespace exports from the original module symbol

## Corresponding TypeScript behavior

Stable TypeScript 5.9 takes a different route in `node_modules/typescript/lib/typescript.js`:

- `getExternalModuleMember`
- `resolveExternalModuleSymbol`
- `getCommonJsExportEquals`

Key point:

- TypeScript's `resolveExternalModuleSymbol` calls `getCommonJsExportEquals(...)`
- `getCommonJsExportEquals(...)` merges non-`export=` members from the original module symbol onto the resolved CommonJS export target

So when TypeScript later does:

- `getExportOfModule(targetSymbol, "Options", ...)`

the merged target symbol already contains `Options`, and no `TS2305` is reported.

## Most likely root cause

The bug was caused by **older `tsgo` not modeling TypeScript's CommonJS `export =` merge behavior closely enough**.

More concretely:

- the package's type alias `Options` lived on the original module symbol
- `resolveExternalModuleSymbol` returned only the resolved `export =` target
- named-export lookup then queried the resolved target instead of a merged/original export container
- therefore `Options` disappeared during lookup

`skipLibCheck` was incidental: it suppressed the package's own declaration-file errors so the incorrect `TS2305` became the visible symptom.

## Surgical fix

If this regression reappears, the obvious fix is to preserve TypeScript's CommonJS merging semantics for `export =` modules when resolving named exports.

In practice, the current local fix already does that by:

- looking up named exports on the original module symbol for `export =` cases
- preserving supplemental type/namespace exports in `getExportsOfModuleWorker`
