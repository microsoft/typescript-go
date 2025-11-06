# Investigation: Issue #1943 - Export Star Re-export with Package.json Exports/Imports

## Issue Summary

When using `export * from ...` to re-export symbols in combination with:
- package.json `exports` subpaths
- package.json `imports` (import maps) with wildcard patterns
- CommonJS module system (`"type": "commonjs"`)
- Declaration files (.d.ts) instead of source files

The re-exported symbols may fail to be resolved with error: "Module has no exported member 'stubDynamicConfig'." This behavior differs from `tsc` 5.9.3 and the native TypeScript LSP.

## Reproduction

### Test Case 1: Source Files (Working)
**Location:** `testdata/tests/cases/compiler/issue1943_export_star_reexport.ts`

This test uses source .ts files directly and appears to work correctly:
- `pkg-exporter` defines source files in `src/`
- `testing.ts` uses `export * from "#pkg-exporter/dep.ts"`
- The import `{ stubDynamicConfig } from "pkg-exporter/testing"` resolves successfully
- Type information is correct: `stubDynamicConfig : () => string`

**Observation:** When source files are available and the module resolution can find them, the export star re-export works as expected.

### Test Case 2: Declaration Files Only (Failing)
**Location:** `testdata/tests/cases/compiler/issue1943_export_star_reexport_dts.ts`

This test uses only .d.ts declaration files (simulating a built/published package) and demonstrates the issue:

```json
{
  "name": "pkg-exporter",
  "type": "commonjs",
  "exports": {
    "./testing": {
      "types": "./dist/testing.d.ts",
      "default": "./dist/testing.js"
    }
  },
  "imports": {
    "#pkg-exporter/*.ts": {
      "types": "./dist/*.d.ts",
      "default": "./dist/*.js"
    }
  }
}
```

**Errors observed:**
1. `/index.ts(2,35): error TS2307: Cannot find module 'pkg-exporter/testing' or its corresponding type declarations.`
2. `/node_modules/pkg-exporter/dist/testing.d.ts(2,15): error TS2307: Cannot find module '#pkg-exporter/dep.ts' or its corresponding type declarations.`

**Root cause hypothesis:** When resolving the `imports` pattern `#pkg-exporter/*.ts` from within a .d.ts file, the module resolution fails to properly apply the pattern matching or resolve to the correct .d.ts file.

## Analysis of Code Paths

### 1. Binder: Export Declaration Binding
**File:** `internal/binder/binder.go`
**Function:** `bindExportDeclaration`

This function collects `ExportStar` symbols from `export *` declarations and adds them to the symbol table. For our test case, this should create an export star entry pointing to the target module.

### 2. Checker: Export Enumeration
**File:** `internal/checker/services.go`
**Function:** `GetExportsOfModule`, `ForEachExportAndPropertyOfModule`

These functions enumerate all exports from a module, including re-exports from `export *` declarations. The checker should:
1. Find the module referenced by the export star
2. Get all its exports
3. Merge them into the current module's export list

**Potential issue:** If the target module cannot be resolved (due to the `imports` pattern not working in .d.ts context), the export star will have no exports to re-export.

### 3. Module Resolution: Exports/Imports Patterns
**Relevant files:**
- `internal/module/*` - Module resolution logic
- `internal/modulespecifiers/*` - Module specifier handling

**Key functions:**
- `loadModuleFromSelfNameReference` - Handles package self-references
- `loadModuleFromExportsOrImports` - Resolves through package.json exports/imports
- `tryGetModuleNameFromExportsOrImports` - Pattern matching for wildcards

**Hypothesis:** When resolving `#pkg-exporter/dep.ts` from within `/node_modules/pkg-exporter/dist/testing.d.ts`:
1. The resolver recognizes it as an `imports` pattern
2. It attempts to match against `#pkg-exporter/*.ts`
3. The pattern matching extracts `dep` as the wildcard match
4. It tries to resolve to `./dist/dep.d.ts` (for types) or `./dist/dep.js` (for default)
5. **Possible failure point:** The pattern matching or file resolution fails, possibly due to:
   - Incorrect file extension handling (`.ts` in pattern vs `.d.ts` or `.js` on disk)
   - Not considering the current file's location when resolving relative to package root
   - Missing logic to resolve from declaration file context

### 4. CommonJS Emit: Export Star Helper
**File:** `internal/printer/helpers.go`
**Function:** `__exportStar` helper emission

When emitting CommonJS, export star declarations are transformed to use the `__exportStar` runtime helper. This is working correctly in the source file test case.

### 5. External Module Info Collection
**File:** `internal/transformers/moduletransforms/*`

The `externalModuleInfoCollector` determines if a module has export stars that need runtime helpers. The flag `hasExportStarsToExportValues` affects emit.

## Hypotheses

### Primary Hypothesis: Import Pattern Resolution in Declaration Files
When resolving an `imports` pattern from within a .d.ts file:
1. The pattern `#pkg-exporter/*.ts` expects files with `.ts` extension
2. But in the built package, only `.d.ts` and `.js` files exist
3. The resolution logic may not correctly map `.ts` extension in the pattern to `.d.ts` for types or `.js` for default
4. This causes the `export * from "#pkg-exporter/dep.ts"` to fail resolution

**Supporting evidence:**
- Test case 1 (source files) works because the actual `.ts` files exist
- Test case 2 (declaration files only) fails because no `.ts` files exist

### Secondary Hypothesis: Package Exports Path Resolution
The error "Cannot find module 'pkg-exporter/testing'" suggests that even the initial `exports` mapping may not be working correctly. This could be due to:
1. The `exports` field not being recognized or processed
2. The types/default condition mapping not being applied correctly
3. Missing logic to resolve subpath patterns in packages

### Tertiary Hypothesis: Symbol Table Merging
Even if resolution succeeds, the symbol table might not correctly merge:
1. Value exports vs type-only exports when only .d.ts available
2. Re-exported symbols when the source came via pattern subpath
3. `SymbolFlags.Function` might be missing in re-export context

## Proposed Diagnostic Steps

1. **Add tracing to module resolution:**
   - Log when `loadModuleFromExportsOrImports` is called
   - Log pattern matching results for `#pkg-exporter/*.ts`
   - Log the resolved file path and whether it exists

2. **Compare symbol tables:**
   - Dump symbol table for `dep.ts` / `dep.d.ts` module
   - Dump symbol table for `testing.ts` / `testing.d.ts` module after export star
   - Compare with `tsc` output to see what symbols are present/missing
   - Check for presence of `SymbolFlags.Function` on `stubDynamicConfig`

3. **Test extension handling:**
   - Modify the pattern to use `#pkg-exporter/*.d.ts` instead
   - Check if resolution succeeds
   - This would confirm if the issue is extension mapping

4. **Simplify reproduction:**
   - Try without wildcards: `"#pkg-exporter/dep": "./dist/dep.d.ts"`
   - If this works, the issue is specifically with wildcard pattern matching

5. **Compare with TypeScript 5.9.3:**
   - Run the same test case through `tsc` 5.9.3
   - Examine how it resolves the pattern
   - Check TypeScript compiler source for relevant resolution logic

## Next Steps

### For Fix Implementation:
1. Identify where imports pattern resolution happens for .d.ts files
2. Add logic to map file extensions in patterns (.ts → .d.ts for types, .ts → .js for default)
3. Ensure package.json `exports` subpaths are correctly resolved
4. Add test coverage for both scenarios

### For Further Investigation:
1. Review TypeScript's implementation of package.json imports/exports resolution
2. Check if there are existing issues or PRs in TypeScript repo about similar problems
3. Test with more complex patterns (nested wildcards, multiple conditions)
4. Verify behavior with ESM vs CommonJS

## Current Status

- [x] Minimal reproduction created in compiler test suite
- [x] Issue confirmed: Declaration-file-only scenario fails
- [x] Baselines accepted showing the error
- [ ] Root cause identified
- [ ] Fix implemented
- [ ] Additional test coverage added

## Related Code References

- Binder: `internal/binder/binder.go` - `bindExportDeclaration`
- Checker: `internal/checker/services.go` - `GetExportsOfModule`, `ForEachExportAndPropertyOfModule`
- Module Resolution: `internal/module/*`, `internal/modulespecifiers/*`
- Emit: `internal/printer/helpers.go` - `__exportStar`
- Transforms: `internal/transformers/moduletransforms/*` - CommonJS module transformer

## Conclusion

The issue is reproducible in the test suite, specifically when:
1. Using package.json `imports` patterns with wildcards
2. Using only declaration files (.d.ts), not source files
3. Attempting to resolve through the pattern from within a .d.ts file

The root cause appears to be related to how import patterns with file extensions are resolved when only built outputs exist. Further investigation is needed to pinpoint the exact code path that needs modification.
