# Fix Summary: Crash with Malformed package.json Imports

## Issue
The TypeScript-Go language server crashed with a panic when processing imports with malformed package.json imports mappings:
```
panic handling request textDocument/diagnostic: should be able to extract TS extension from string that passes IsDeclarationFileName
```

## Investigation Results

### What causes the crash?
1. A package.json has a malformed imports pattern mapping: `"./src/*ts"` (missing dot before `ts`)
2. An import specifier like `"#/b."` is used (ending with dot, no extension)
3. The pattern substitution resolves to a valid file (`/src/b.ts`)
4. The `resolvedUsingTsExtension` flag is set to true (because package.json value doesn't end with `.ts`)
5. The checker tries to extract the TS extension from `"#/b."` for error reporting
6. `TryExtractTSExtension("#/b.")` returns empty string (because it doesn't end with a valid TS extension)
7. **PANIC!** - Code assumed extraction would always succeed when `resolvedUsingTsExtension` is true

### Root Cause
The code in `/internal/checker/checker.go` at lines 14817 and 14833 incorrectly assumed that:
- If `resolvedUsingTsExtension` is true, then
- The `moduleReference` must have a TS extension that can be extracted

This assumption is **invalid** when:
- Package.json pattern mappings are used
- The mapping is malformed
- The import specifier doesn't actually contain a TS extension

## Solution
Modified `/internal/checker/checker.go` to gracefully handle the case where `TryExtractTSExtension` returns an empty string. Instead of panicking, the code now skips the error reporting for that specific case. This is correct behavior because:
- The specific errors about TS extensions don't apply when the module reference doesn't have a TS extension
- The resolution was successful (the file was found)
- Other diagnostics (if any) about the malformed package.json may still be reported

## Changes Made
1. **`internal/checker/checker.go`**: Replaced panic with conditional error reporting
   - Line 14809-14824: Skip error reporting if TS extension extraction fails
   - Line 14825-14840: Skip error reporting if TS extension extraction fails
   
2. **`internal/fourslash/tests/manual/packageJsonImportsMalformed_test.go`**: Added test case
   - Reproduces the exact crash scenario
   - Verifies the fix prevents the panic
   - Tests with malformed package.json pattern: `"./src/*ts"` instead of `"./src/*.ts"`

3. **`INVESTIGATION.md`**: Added detailed investigation notes
   - Documents the root cause analysis
   - Explains why the assumption was violated
   - Provides context for future maintainers

## Testing
- ✅ New test `TestPackageJsonImportsMalformed` passes (reproduces the original crash)
- ✅ All existing tests in `internal/fourslash/tests/manual/` pass
- ✅ No regressions detected
- ✅ CodeQL security check passed

## Impact
- **Positive**: Prevents crash in edge case with malformed package.json
- **Neutral**: No change to normal behavior - existing functionality preserved
- **No Breaking Changes**: The fix only affects error handling in an edge case

## Security Summary
No security vulnerabilities introduced. The fix makes the code more robust by handling an edge case that previously caused a crash.
