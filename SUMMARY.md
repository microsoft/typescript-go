# Investigation Summary: Invalid character error on backslash in JSDoc comment

## Problem
tsgo incorrectly reports error TS1127 "Invalid character" when encountering a backslash in a JSDoc comment in JavaScript files, while TypeScript 5.9 does not report this error.

### Reproduction
```js
/**
 * \
 */
function foo() {}
```

**Expected (TypeScript 5.9):** No error  
**Actual (tsgo before fix):** `error TS1127: Invalid character.`

## Root Cause
The issue was in `internal/scanner/scanner.go` in the `ScanJSDocToken()` function (lines 1359-1371).

When the scanner encountered a backslash character in JSDoc context, it checked if it was followed by a valid Unicode escape sequence. If not, it called `scanInvalidCharacter()` which produced an error diagnostic.

This was incorrect because:
1. Backslashes should be allowed as regular characters in JSDoc comments
2. TypeScript 5.9 allows backslashes in JSDoc comments without errors
3. The scanner's `ScanJSDocCommentTextToken()` function already handles backslashes correctly

## Solution
Modified `ScanJSDocToken()` to return `KindUnknown` without calling `scanInvalidCharacter()` when a backslash is not followed by a valid Unicode escape. This allows the parser to handle the backslash based on its state machine, matching TypeScript 5.9 behavior.

**Code change:**
```go
// Before:
} else {
    s.scanInvalidCharacter()  // Reports error
}

// After:
} else {
    // Backslash not followed by valid unicode escape - treat as unknown token
    // without reporting an error. The parser will handle this based on its state.
    s.pos++
    s.token = ast.KindUnknown
}
```

## Testing
1. Created test case: `testdata/tests/cases/compiler/jsdocBackslash.ts`
2. Verified fix with multiple backslash patterns in JSDoc comments
3. Confirmed behavior matches TypeScript 5.9 (no error reported)
4. Scanner still returns `KindUnknown` for backslashes but doesn't produce error diagnostics

## Files Changed
- `internal/scanner/scanner.go`: Fixed the backslash handling in `ScanJSDocToken()`
- `testdata/tests/cases/compiler/jsdocBackslash.ts`: Added test case
- `INVESTIGATION.md`: Detailed investigation notes

## Verification
```bash
# Both commands now produce no errors:
./built/local/tsgo --allowJs --noEmit testdata/tests/cases/compiler/jsdocBackslash.ts
npx tsc --allowJs --noEmit testdata/tests/cases/compiler/jsdocBackslash.ts
```

## Additional Context
The investigation revealed:
- JSDoc diagnostics in JS files are stored separately in `p.jsdocDiagnostics`
- Error callback chain: `scanner.errorAt()` → `scanner.onError()` → `parser.scanError()` → `parser.parseErrorAtRange()`
- The JSDoc parser state machine determines when to use `ScanJSDocToken()` vs `ScanJSDocCommentTextToken()`
- Backslashes are correctly handled in `ScanJSDocCommentTextToken()` but were incorrectly flagged as errors in `ScanJSDocToken()`
