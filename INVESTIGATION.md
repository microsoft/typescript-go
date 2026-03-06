# Investigation: Fourslash Tests Failing for Unconfigured .js Files

## Problem Summary

Fourslash tests were failing when making requests to .js files without explicit `@allowJs: true` or `@checkJs: true` compiler options.

**Error:** `no project found for URI file:///d.js`

## Root Cause

The issue was in `internal/project/project.go` in the `NewInferredProject` function. The inferred project is designed to support JavaScript files by default (by setting `AllowJs: true` in its default compiler options). However, when compiler options were explicitly provided (even without setting AllowJs), the default AllowJs setting was bypassed.

### Code Flow

1. Fourslash tests call `SetCompilerOptionsForInferredProjects` with compiler options that don't include AllowJs
2. When the inferred project is created via `NewInferredProject`, it checks if `compilerOptions == nil`
3. If nil, it uses defaults including `AllowJs: core.TSTrue`
4. If non-nil (which was the case), it used the provided options WITHOUT the AllowJs default
5. Without AllowJs, .js files were not properly included in the inferred project
6. When a language service request was made for a .js file, no project could be found

## Investigation Process

1. **Reproduced the issue** - Created a minimal test case that demonstrates the problem
2. **Traced the code flow** - Followed the execution path from DidOpen through project creation
3. **Identified the root cause** - Found that AllowJs default was not being applied when compilerOptions was non-nil
4. **Implemented and tested the fix** - Added logic to apply AllowJs default when TSUnknown
5. **Verified the fix** - Ran all tests to ensure no regressions

## Test Cases

Created `internal/fourslash/tests/test_js_file_issue_test.go` with two tests:
1. `TestJsFileCompletions` - with `@checkJs: true` (works before and after fix)
2. `TestJsFileWithoutCheckJs` - without any JS option (fails before fix, passes after fix)

Both tests use multiple .js files to ensure proper project setup and file inclusion.

## Fix

Modified `NewInferredProject` in `internal/project/project.go` to apply the AllowJs default even when compiler options are provided, if AllowJs is not explicitly set:

```go
} else {
    // Apply inferred project defaults for options that are not explicitly set
    if compilerOptions.AllowJs == core.TSUnknown {
        compilerOptions.AllowJs = core.TSTrue
    }
}
```

This ensures that the inferred project always supports JavaScript files by default, which is consistent with the intended behavior and matches TypeScript's language server behavior.

## Why This Matters

The inferred project is used when files are opened without a tsconfig.json. It should be permissive by default to provide a good developer experience. JavaScript files are commonly mixed with TypeScript files in modern projects, so the inferred project must support them without requiring explicit configuration.

## Verification

- Both test cases now pass
- All existing project tests continue to pass (2.3s runtime)
- All fourslash tests continue to pass
- The fix is minimal and focused on the root cause
- Code review found no issues
- No security vulnerabilities detected

## Future Considerations

Other compiler options in the default set (like `AllowNonTsExtensions`) might benefit from similar treatment, but AllowJs is the most critical for file inclusion. If similar issues arise with other options, they can be addressed separately with the same pattern.

