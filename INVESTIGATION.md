# Investigation: Invalid character error on backslash in JSDoc comment in JS files

## Issue Description
When tsgo processes a JavaScript file with a backslash in a JSDoc comment, it produces error TS1127 "Invalid character", while TypeScript 5.9 does not report this error.

### Reproduction
File: `index.js`
```js
/**
 * \
 */
```

Expected: No error (as with TypeScript 5.9)
Actual with tsgo: `index.js:2:4 - error TS1127: Invalid character.`

## Root Cause Analysis

### Investigation Steps

1. **Identified error source**: The error TS1127 "Invalid character" is defined in `internal/diagnostics/diagnostics_generated.go:193`.

2. **Found error trigger**: The error is produced by `scanInvalidCharacter()` in `internal/scanner/scanner.go:1987-1992`.

3. **Traced scanner execution**: Using a test scanner, confirmed that when scanning the JSDoc comment `/** \n * \ \n */`:
   - The backslash at position 7 produces a `KindUnknown` token
   - The scanner calls `scanInvalidCharacter()` which reports the error

4. **Located problematic code**: In `internal/scanner/scanner.go`, the `ScanJSDocToken()` function at lines 1359-1368:
   ```go
   case '\\':
       s.pos--
       cp := s.peekUnicodeEscape()
       if cp >= 0 && IsIdentifierStart(cp) {
           s.tokenValue = string(s.scanUnicodeEscape(true)) + s.scanIdentifierParts()
           s.token = GetIdentifierToken(s.tokenValue)
       } else {
           s.scanInvalidCharacter()  // ← THIS IS THE PROBLEM
       }
       return s.token
   ```

5. **Understood the context**: The JSDoc parser state machine in `internal/parser/jsdoc.go:156-270` shows:
   - After `/**` and newline, state is `jsdocStateBeginningOfLine`
   - After whitespace ` `, state remains the same
   - After asterisk `*`, state becomes `jsdocStateSawAsterisk`
   - After whitespace ` `, state remains `jsdocStateSawAsterisk`
   - When backslash `\` is encountered, since state is NOT `jsdocStateSavingComments`, it calls `nextTokenJSDoc()` instead of `nextJSDocCommentTextToken()`
   - `nextTokenJSDoc()` calls `ScanJSDocToken()` which treats the backslash as invalid

### Why This Is Wrong

A backslash in a JSDoc comment should be allowed as regular comment text. The scanner function `ScanJSDocCommentTextToken()` correctly handles backslashes by including them in the comment text token (line 1277). However, `ScanJSDocToken()` is being called instead in certain parsing states, and it incorrectly treats a standalone backslash as an "invalid character."

In TypeScript 5.9, a backslash in JSDoc comment text is allowed and does not produce an error.

### The Fix

The issue is in `internal/scanner/scanner.go` in the `ScanJSDocToken()` function (lines 1359-1368).

When `ScanJSDocToken()` encounters a backslash that is not followed by a valid Unicode escape sequence, it should NOT call `scanInvalidCharacter()`. Instead, it should treat the backslash as an unknown token WITHOUT reporting an error, similar to how the default case (line 1389-1392) handles unrecognized characters.

**Current code:**
```go
case '\\':
    s.pos--
    cp := s.peekUnicodeEscape()
    if cp >= 0 && IsIdentifierStart(cp) {
        s.tokenValue = string(s.scanUnicodeEscape(true)) + s.scanIdentifierParts()
        s.token = GetIdentifierToken(s.tokenValue)
    } else {
        s.scanInvalidCharacter()  // ERROR: This reports an error
    }
    return s.token
```

**Proposed fix:**
```go
case '\\':
    s.pos--
    cp := s.peekUnicodeEscape()
    if cp >= 0 && IsIdentifierStart(cp) {
        s.tokenValue = string(s.scanUnicodeEscape(true)) + s.scanIdentifierParts()
        s.token = GetIdentifierToken(s.tokenValue)
    } else {
        // Backslash not followed by valid unicode escape - treat as unknown token
        // without reporting an error (will be handled as comment text)
        s.pos++
        s.token = ast.KindUnknown
    }
    return s.token
```

## Test Case

Created test file: `testdata/tests/cases/compiler/jsdocBackslash.ts`

```typescript
// @allowJs: true
// @noEmit: true
// @filename: test.js

/**
 * \
 */
function foo() {}
```

This test should pass without errors, matching TypeScript 5.9 behavior.

## Additional Notes

- JSDoc diagnostics in JS files are stored separately in `p.jsdocDiagnostics` (see `internal/parser/jsdoc.go:134-138`)
- The error callback chain: `scanner.errorAt()` → `scanner.onError()` → `parser.scanError()` → `parser.parseErrorAtRange()`
- Backslashes ARE correctly handled in `ScanJSDocCommentTextToken()` (line 1258-1284), but that function is only called when the parser state is `jsdocStateSavingComments`

## Recommendation

Apply the proposed fix to `internal/scanner/scanner.go` at lines 1359-1368 to allow backslashes in JSDoc comments without reporting an error.
