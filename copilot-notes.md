# Investigation: JSDoc Type Assertion Panic in getTokenAtPosition

## Problem Statement

The `getTokenAtPosition` function panics with the message:
```
panic: did not expect KindParenthesizedExpression to have KindIdentifier in its trivia
```

When processing code like:
```javascript
function foo(x) {
    const s = /**@type {string}*/(x)  // Go-to-definition on x causes panic
}
```

## Initial (Incorrect) Fix

My initial fix simply removed the panic and returned `current` when an identifier was found during trivia scanning. However, @andrewbranch correctly pointed out that:
- Identifiers are NEVER allowed in the trivia of nodes
- JSDoc nodes should be fully-formed subtrees with their identifiers accounted for in the AST
- The panic is there for a reason and should not be removed

## Investigation Plan

1. Understand how JSDoc nodes are structured in the AST
2. Compare working JSDoc cases (that don't panic) with the failing JSDoc type assertion case
3. Identify why the identifier is appearing in trivia for type assertions but not other JSDoc constructs
4. Find the root cause and proper fix

## Step 1: Understanding the AST Structure

From the test output, here's what I observed:

### Case 1: JSDoc Type Assertion (The Failing Case)
```
KindParenthesizedExpression [31..54) Flags=0x210000
  JSDoc: KindJSDoc [31..51)
  KindAsExpression [52..53)
    KindIdentifier [52..53)
    KindStringKeyword [42..48)
```

Key observations:
- The `ParenthesizedExpression` has the `HasJSDoc` flag (0x200000 in Flags=0x210000)
- The JSDoc is at [31..51), which includes the `/**@type {string}*/` part
- The JSDoc is attached as metadata to the `ParenthesizedExpression`
- The `AsExpression` starts at position 52 (the identifier `x`)
- The `StringKeyword` at [42..48) is from inside the JSDoc comment

### Case 2: Regular JSDoc Comment (Works Fine)
```
KindFunctionDeclaration [0..62) Flags=0x210000
  JSDoc: KindJSDoc [0..28)
  KindIdentifier [37..41)
  KindParameter [42..43)
    KindIdentifier [42..43)
    KindStringKeyword [15..21)
```

Key observations:
- The JSDoc is properly attached to the `FunctionDeclaration`
- The `StringKeyword` at [15..21) is also from the JSDoc
- The parameter identifier at [42..43) is a separate node

## Step 2: Identifying the Problem

**ROOT CAUSE FOUND!**

When `getTokenAtPosition` is called with position 52 (the identifier `x`):

1. It starts at `ParenthesizedExpression` [31..54)
2. Calls `VisitEachChildAndJSDoc`, which:
   - First visits JSDoc children if `HasJSDoc` flag is set (the JSDoc is at [31..51), doesn't contain position 52)
   - Then visits regular children via `VisitEachChild` (the `AsExpression` at [52..53))
3. The `AsExpression` IS visited, and `testNode` would return 0 (match)
4. **BUT** - the `visitNode` function has this check:
   ```go
   if node != nil && node.Flags&ast.NodeFlagsReparsed == 0 && next == nil {
   ```
5. The `AsExpression` has Flags=0x10008, which includes `NodeFlagsReparsed` (0x8)
6. Because the Reparsed flag is set, `visitNode` **skips the AsExpression entirely**!
7. Since no child is found, the function falls back to using the scanner
8. The scanner scans from `left` (which is after the JSDoc) and encounters the identifier `x`
9. The code panics because it doesn't expect an identifier in trivia

## Step 3: Understanding "Reparsed" Nodes

The `NodeFlagsReparsed` flag is used to mark nodes that are created from reparsing JSDoc comments. When a JSDoc type annotation like `/**@type {string}*/` is encountered, it gets "reparsed" into proper AST nodes (in this case, an `AsExpression`).

The problem is that `getTokenAtPosition` explicitly skips reparsed nodes during traversal:
```go
if node != nil && node.Flags&ast.NodeFlagsReparsed == 0 && next == nil {
```

This is intentional - reparsed nodes represent synthetic AST nodes created from JSDoc, and their positions overlap with the JSDoc comment text. The code is designed to skip them to avoid confusion.

However, in the case of JSDoc type assertions, the `AsExpression` is the ONLY child of the `ParenthesizedExpression` (besides the JSDoc itself). When it's skipped, there's no other child to navigate to, so the scanner kicks in and finds the identifier.

## Step 4: Comparing with Working JSDoc Cases

Let me check how other JSDoc cases work:

**Case 2: Regular @param tag**
```
KindFunctionDeclaration [0..62)
  JSDoc: KindJSDoc [0..28)
  KindParameter [42..43)
    KindIdentifier [42..43)
    KindStringKeyword [15..21)  <-- Reparsed node, but it's inside the Parameter
```

In this case, the `StringKeyword` (the type from JSDoc) is a child of the `Parameter` node. When we navigate to the Parameter node, we have a real, non-reparsed identifier at [42..43) that we can find. The reparsed `StringKeyword` is a sibling, not the only path to the identifier.

**Case 1: JSDoc type assertion (failing)**
```
KindParenthesizedExpression [31..54)
  JSDoc: KindJSDoc [31..51)
  KindAsExpression [52..53)  <-- Reparsed node, ONLY child
    KindIdentifier [52..53)
```

In this case, the `AsExpression` IS the only child (besides the JSDoc). When it's skipped, there's NO other path to the identifier.

## Step 5: The Solution

The issue is that when navigating into a JSDoc type assertion, we need to be able to navigate through reparsed nodes to reach the actual identifiers within them. The current code correctly skips reparsed nodes to avoid returning positions within JSDoc comments, but it doesn't handle the case where the reparsed node is the ONLY way to reach a real token.

Looking at `visitNodeList`, there's already special handling for reparsed nodes:
```go
if match && nodes[index].Flags&ast.NodeFlagsReparsed != 0 {
    // filter and search again
    nodes = core.Filter(nodes, func(node *ast.Node) bool {
        return node.Flags&ast.NodeFlagsReparsed == 0
    })
    // ... search again without reparsed nodes
}
```

The problem is that `visitNode` (used for single children) doesn't have this fallback logic. It just skips reparsed nodes unconditionally.

**The Fix:** Modify `visitNode` to match the behavior of `visitNodeList`. When a reparsed node is encountered that would match the position, we should still process it if there are no non-reparsed alternatives. This allows navigation through JSDoc type assertions while maintaining the invariant that identifiers never appear in trivia.

Specifically, the fix is to:
1. When `visitNode` encounters a reparsed node that matches, don't immediately skip it - instead, save it in a `reparsedNext` variable
2. After visiting all children, if no non-reparsed match was found (`next == nil`), use the reparsed match
3. Reset `reparsedNext` at the end of each iteration

This way:
- Normal cases continue to prefer non-reparsed nodes (e.g., when a parameter has both a JSDoc type and a real identifier)
- JSDoc type assertion cases can navigate through the reparsed AsExpression to reach the identifier
- The invariant is maintained: we never return an identifier found in trivia

## Implementation

The changes made to `internal/astnav/tokens.go`:

1. Added `reparsedNext` variable to track reparsed node matches as a fallback
2. Modified `visitNode` to:
   - Still test reparsed nodes
   - Save matching reparsed nodes in `reparsedNext` instead of skipping them entirely
   - Prefer non-reparsed nodes by only setting `next` for non-reparsed matches
3. Added logic after `VisitEachChildAndJSDoc` to use `reparsedNext` if `next` is nil
4. Reset both `next` and `reparsedNext` at the end of each loop iteration

This maintains backward compatibility while fixing the JSDoc type assertion panic. The panic condition is kept intact - identifiers should never appear in trivia.

## Testing

All existing tests pass, including:
- The new test cases for JSDoc type assertions
- All baseline tests for `GetTokenAtPosition` and `GetTouchingPropertyName`
- All other astnav tests

The fix correctly handles:
- JSDoc type assertions like `/**@type {string}*/(x)`
- Regular JSDoc comments (unchanged behavior)
- All other token position lookups (unchanged behavior)
