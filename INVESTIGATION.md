# Investigation: Signature Help Crash with Binding Pattern Parameters

## Problem Statement
When requesting signature help for a function with binding pattern parameters (e.g., `function foo({})`), the language server crashes with:
```
panic handling request textDocument/signatureHelp: Unhandled case in Node.Text: *ast.BindingPattern
```

## Root Cause Analysis

### Call Stack
1. `ls.(*LanguageService).createSignatureHelpParameterForParameter()` in `signaturehelp.go:594`
   - Calls `getDocumentationFromDeclaration()` to get JSDoc for parameter
2. `ls.(*LanguageService).getDocumentationFromDeclaration()` in `hover.go:80`
   - Calls `getJSDocOrTag()` to find JSDoc comments
3. `getJSDocOrTag()` in `hover.go:455`
   - For parameters, tries to match JSDoc @param tags by calling `node.Name().Text()`
4. `ast.(*Node).Text()` in `ast.go:345`
   - **Panics because it doesn't handle `BindingPattern` nodes**

### Type System Analysis

From `internal/ast/ast.go`:
- Line 2100: `BindingName = Node // Identifier | BindingPattern`
- Line 2089: `DeclarationName = Node // ... | BindingPattern | ...`

A parameter's name can be either:
- An `Identifier` (e.g., `function foo(x)`)
- A `BindingPattern` (e.g., `function foo({x})` or `function foo([x])`)

The `Text()` method at line 310-346 only handles:
- Identifiers
- Literals
- JSDoc nodes
- A few other node types

**It does NOT handle `BindingPattern` nodes**, causing the panic.

### TypeScript Behavior

Looking at `_submodules/TypeScript/src/compiler/utilities.ts`:

```typescript
export function getParameterSymbolFromJSDoc(node: JSDocParameterTag): Symbol | undefined {
    if (node.symbol) {
        return node.symbol;
    }
    if (!isIdentifier(node.name)) {  // <-- Key check!
        return undefined;
    }
    const name = node.name.escapedText;
    const parameter = find(decl.parameters, p => 
        p.name.kind === SyntaxKind.Identifier &&  // <-- Only matches identifiers
        p.name.escapedText === name
    );
    return parameter && parameter.symbol;
}
```

TypeScript explicitly:
1. Checks if the parameter name is an identifier before trying to match JSDoc
2. Skips binding patterns when matching JSDoc @param tags

This makes sense because JSDoc @param syntax doesn't support binding patterns:
```typescript
/** @param {a, b} */ // Not valid JSDoc syntax
function foo({a, b}) {}
```

## Solution

Added a check in `internal/ls/hover.go` at line 454-462 to skip JSDoc matching for parameters with binding patterns:

```go
case ast.IsParameter(node):
    // Parameters with binding patterns (e.g., {a, b}) don't have a simple name to match against JSDoc.
    // TypeScript also skips JSDoc matching for binding patterns.
    name := node.Name()
    if ast.IsBindingPattern(name) {
        return nil
    }
    return getMatchingJSDocTag(c, node.Parent, name.Text(), isMatchingParameterTag)
```

This fix:
1. Prevents the crash by avoiding the call to `Text()` on binding patterns
2. Aligns with TypeScript's behavior
3. Is semantically correct (JSDoc doesn't support binding patterns anyway)

## Test Coverage

Added three test cases in `internal/fourslash/tests/signatureHelpBindingPattern_test.go`:

1. **TestSignatureHelpBindingPattern**: Basic object binding pattern `{}`
2. **TestSignatureHelpBindingPatternWithJSDoc**: Object binding pattern with JSDoc on function
3. **TestSignatureHelpArrayBindingPattern**: Array binding pattern `[x, y]`

All tests verify that signature help works without crashing.

## Alternative Solutions Considered

### Option 1: Implement `Text()` for BindingPattern
Could add a case in `Node.Text()` to handle binding patterns by returning some representation like `"{}"` or `"[...]"`.

**Rejected because**:
- Would need to decide on a format (what text to return?)
- Doesn't solve the semantic issue (JSDoc doesn't support binding patterns)
- The real problem is trying to match JSDoc on binding patterns, not the missing `Text()` implementation

### Option 2: Handle at the higher level in signaturehelp.go
Could check for binding patterns in `createSignatureHelpParameterForParameter()` before calling `getDocumentationFromDeclaration()`.

**Rejected because**:
- `getJSDocOrTag()` is more general-purpose and used in other places
- The fix at the `getJSDocOrTag()` level is more defensive and prevents similar issues elsewhere
- Matches TypeScript's architecture where the check is in the JSDoc matching logic

## Impact

- **Fixed**: Crash when requesting signature help on functions with binding pattern parameters
- **No breaking changes**: Signature help still works, just without (unsupported) JSDoc for binding patterns
- **Test coverage**: Added comprehensive tests for both object and array binding patterns
- **Alignment with TypeScript**: Behavior now matches TypeScript's implementation

## Files Modified

1. `internal/ls/hover.go` - Added binding pattern check in `getJSDocOrTag()`
2. `internal/fourslash/tests/signatureHelpBindingPattern_test.go` - Added test cases
