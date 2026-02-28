# Investigation: JSX Parsing Fails with Ternary + Nested JSX Attributes + Object Literals

## Summary

JSX parsing produces spurious syntax errors in ternary expressions when:
1. The truthy branch has a parenthesized expression with an identifier on a separate line (e.g., `(\n  x\n)`)
2. The falsy branch contains JSX with multi-line attributes
3. The JSX attributes contain nested JSX with function calls that include multi-property object literals

## Minimal Reproduction

```tsx
// @jsx: preserve
const a = (
  <div>
    {true ? (
      x
    ) : (
      <div
        label={
          <div>{f(1, {a:1,b:2})}</div>
        }
      ></div>
    )}
  </div>
)
```

This produces ~14 syntax errors starting with `',' expected` at the `label` attribute.

## Root Cause

The bug is in `parseAssignmentExpressionOrHigherWorker` → `tryParseParenthesizedArrowFunctionExpression` in `internal/parser/parser.go`.

### Detailed Trace

When parsing the ternary expression `true ? (\n x\n) : (\n <div...`, the parser:

1. **Parses `true`** as the condition
2. **Sees `?`**, enters `parseConditionalExpressionRest`
3. **Parses the truthy branch** via `parseAssignmentExpressionOrHigherWorker`
4. The truthy branch starts with `(`. `isParenthesizedArrowFunctionExpression()` returns `TSUnknown` (because `(x)` matches the pattern for potential arrow function `(param) => ...`)
5. **Speculative parse begins** via `tryParseParenthesizedArrowFunctionExpression` → `parsePossibleParenthesizedArrowFunctionExpression` → `parseParenthesizedArrowFunctionExpression(false /*allowAmbiguity*/, ...)`
6. The speculative parser:
   - Parses `(x)` as a parameter list with one parameter `x`
   - After `)`, the token is `:` — but this is the **ternary colon**, not a return type colon!
   - `hasReturnColon = true` (line 4309 in parser.go)
   - `parseReturnType(KindColonToken, false)` consumes the `:` and tries to parse a type

7. **Type parsing gone wrong**: Starting from `(`, the type parser:
   - Enters `parseParenthesizedType()`
   - Inside the parens, sees `<` → treats it as function type with type params
   - Tries to parse `<div\n label={\n <div>{f(1,...` as a type with malformed type parameters
   - Through various error recovery paths, consumes tokens up to `{` (the `{` of `{a:1,b:2}`)

8. **False positive on `{`**: At line 4325:
   ```go
   if !allowAmbiguity && p.token != ast.KindEqualsGreaterThanToken && p.token != ast.KindOpenBraceToken {
       return nil
   }
   ```
   Since `p.token == KindOpenBraceToken`, this check passes — the parser mistakenly thinks `{` could be the start of an arrow function body.

9. **Arrow body parsing**: The parser continues to parse `{a:1,` as an arrow function body, consuming more tokens. After body parsing, `p.token` is `KindColonToken` (the `:` in `b:2`).

10. **Final check also passes** (line 4353-4362):
    ```go
    if !allowReturnTypeInArrowFunction && hasReturnColon {
        if p.token != ast.KindColonToken {
            return nil
        }
    }
    ```
    Since `p.token == KindColonToken`, this check also passes (it's supposed to check for a second `:` in a ternary, and coincidentally finds `:` from `b:2`).

11. **Speculative parse succeeds when it should fail**. The result is non-nil, so `tryParseParenthesizedArrowFunctionExpression` does NOT rewind. The parser has committed to an incorrect arrow function interpretation.

12. **Cascading errors**: Everything after this point is parsed incorrectly, producing the spurious syntax errors.

### Why it works with string literals

When the truthy branch is `("hello")` instead of `(x)`:
- `isParenthesizedArrowFunctionExpression()` returns `TSFalse` because `(` followed by a string literal cannot be a parameter
- The speculative arrow function parse is never attempted
- The ternary parses correctly

### Why it works on a single line

When everything is on one line, the type parameter error recovery inside the speculative type parse consumes different tokens, ultimately landing on a token that ISN'T `{` or `KindColonToken`, causing the speculative parse to correctly fail and rewind.

## Comparison with TypeScript Reference

The TypeScript reference (`_submodules/TypeScript/src/compiler/parser.ts`) has the same speculative parsing logic. The difference appears to be in **type parsing error recovery** behavior during the speculative parse. When the TS reference tries to parse `(<div\n label=...` as a type, its error recovery consumes different tokens than the Go port, resulting in a different final token position that causes the speculative parse to correctly fail.

The Go port is also missing the `hasJSDocFunctionType` check at the equivalent of TS line 5490:
```typescript
const hasJSDocFunctionType = unwrappedType && isJSDocFunctionType(unwrappedType);
if (!allowAmbiguity && token() !== SyntaxKind.EqualsGreaterThanToken && (hasJSDocFunctionType || token() !== SyntaxKind.OpenBraceToken)) {
```
vs Go line 4325:
```go
if !allowAmbiguity && p.token != ast.KindEqualsGreaterThanToken && p.token != ast.KindOpenBraceToken {
```
While this specific difference doesn't directly cause this bug (the unwrapped type isn't a JSDocFunctionType), it should still be ported for correctness.

## Potential Fix Approaches

1. **Fix type parsing error recovery** to match the TypeScript reference behavior more precisely. This would require detailed comparison of how both parsers handle malformed type syntax during speculative parsing.

2. **Add additional guards** in the speculative arrow function parsing to reject cases where the return type parsing consumed JSX-like content. For example, checking if the parsed return type contains error nodes from JSX context.

3. **Port the `hasJSDocFunctionType` check** from the TS reference (low priority for this specific bug but still a correctness fix).

## Key Files

- `internal/parser/parser.go`:
  - `parseAssignmentExpressionOrHigherWorker` (line ~4011)
  - `tryParseParenthesizedArrowFunctionExpression` (line ~4250)
  - `parseParenthesizedArrowFunctionExpression` (line ~4271)
  - `parseConditionalExpressionRest` (line ~4479)
  - `parseJsxExpression` (line ~4787)
  - `parseJsxOpeningOrSelfClosingElementOrOpeningFragment` (line ~4841)

## Test Case

See `testdata/tests/cases/compiler/jsxTernaryWithObjectInAttribute.tsx`
