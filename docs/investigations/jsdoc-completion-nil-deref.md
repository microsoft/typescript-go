# Investigation: Nil Pointer Dereference in Completions After JSDoc Comment

**Date:** 2025-07-11  
**Status:** Root cause identified; fix candidates proposed  
**Failing test:** `TestCompletionsAfterJSDoc` in `internal/fourslash/tests/gen/completionsAfterJSDoc_test.go`

---

## Problem Description

Requesting IDE completions when the cursor is positioned at the **beginning of a class or interface member** that is **preceded by a JSDoc comment** causes a nil pointer dereference (panic) in the completion engine.

### Reproduction

The exact scenario is captured by the existing (skipped) test `TestCompletionsAfterJSDoc`:

```typescript
export interface Foo {
  /** JSDoc */
  /**/foo(): void;   // <-- cursor at start of "foo"
}
```

Expected: keyword completions for interface element context (e.g. `readonly`)  
Actual: panic / nil pointer dereference in `tryGetClassLikeCompletionSymbols`

---

## Affected Files

| File | Relevant Location |
|------|------------------|
| `internal/ls/completions.go` | `tryGetClassLikeCompletionSymbols` closure (~line 1437), nil deref at line 1447 |
| `internal/ls/completions.go` | `tryGetObjectTypeDeclarationCompletionContainer` (~line 4090) |
| `internal/ls/completions.go` | `getRelevantTokens` (~line 2750) |
| `internal/astnav/tokens.go` | `FindPrecedingTokenEx` (~line 290), `findRightmostValidToken` (~line 422) |
| `internal/astnav/tokens.go` | `getNodeVisitor` / `wrappedVisitNodes` (~line 630) |
| `internal/ast/utilities.go` | `IsJSDocSingleCommentNode` (line 2958), `IsJSDocSingleCommentNodeList` (line 2940) |
| `internal/fourslash/tests/gen/completionsAfterJSDoc_test.go` | Existing skipped test |

---

## Root Cause Chain

### Step 1 — Token position analysis

For the test source (cursor at position 40, the start of `foo`):

```
pos 21: '{'   ← interface body open brace
pos 22: '\n'
pos 25: '/'   ← start of /** JSDoc */
pos 37: '\n'  ← JSDoc node's End
pos 40: 'f'   ← cursor position, start of "foo"
```

AST structure of the `MethodSignature`:

```
MethodSignature Pos=22 End=52 (hasJSDoc=true)
  [JSDoc]  Pos=22 End=37
    KindJSDocText Pos=25 End=35
  KindIdentifier (foo) Pos=22 End=43   ← .Pos() == 22 (includes JSDoc range in trivia)
  KindOpenParenToken ...
  KindVoidKeyword ...
```

> Note: The identifier node's `.Pos()` is 22 (beginning of the enclosing member's range,
> including leading JSDoc trivia), while `GetStartOfNode(identifier, file, true)` = 40
> (the actual `f` character, skipping JSDoc + whitespace trivia).

---

### Step 2 — `FindPrecedingToken(file, 40)` returns `nil`

`GetStartOfNode(identifier, file, includeJSDoc=true)` returns 40 (the real start after JSDoc
trivia). Because `40 >= position(40)`, `lookInPreviousChild = true`, which triggers the
"look in preceding JSDoc" path in `FindPrecedingTokenEx`:

```go
// internal/astnav/tokens.go ~line 360
lookInPreviousChild := start >= position || !isValidPrecedingNode(foundChild, sourceFile)
if lookInPreviousChild {
    if position >= foundChild.Pos() {
        // Find jsdoc preceding the foundChild.
        var jsDoc *ast.Node
        // ... finds the JSDoc node at Pos=22 End=37
        if jsDoc != nil {
            if !excludeJSDoc {
                return find(jsDoc)   // ← enters JSDoc recursion
            }
            ...
        }
        ...
    }
}
```

`find(jsDoc)` calls `findRightmostValidToken(37, file, jsDoc, -1, false)`.

Inside `findRightmostValidToken`, `VisitEachChildAndJSDoc` is called on the JSDoc node.
This invokes `getNodeVisitor`'s `wrappedVisitNodes`:

```go
// internal/astnav/tokens.go ~line 645
wrappedVisitNodes = func(n *ast.NodeList, v *ast.NodeVisitor) *ast.NodeList {
    if ast.IsJSDocSingleCommentNodeList(n) {
        return n   // ← SKIP the JSDocText NodeList entirely!
    }
    return visitNodes(n, v)
}
```

`IsJSDocSingleCommentNode` (in `internal/ast/utilities.go` line 2958) returns `true` for
`/** JSDoc */` because the JSDoc node has **exactly one comment node** (a `KindJSDocText`):

```go
// In Strada, if a JSDoc node has a single comment, that comment is represented
// as a string property as a simplification, and therefore that comment is not
// visited by `forEachChild`.
func IsJSDocSingleCommentNode(node *Node) bool {
    return hasComment(node.Kind) && node.CommentList() != nil && len(node.CommentList().Nodes) == 1
}
```

This "Strada optimization" intentionally makes single-comment JSDoc nodes invisible to
child traversal. As a result:

- `hasChildren = false` inside `findRightmostValidToken`
- Case 3 triggers: `n(jsDoc) == containingNode(jsDoc)` → returns `nil`
- `find(jsDoc)` returns `nil`
- Back in `find(methodSignature)`: returns `nil`
- `FindPrecedingToken(file, 40)` = **`nil`**

For comparison: `FindPrecedingTokenEx(file, 40, nil, excludeJSDoc=true)` correctly returns
the `KindOpenBraceToken` (`{`, position 21) because it skips JSDoc entirely and falls back
to `findRightmostValidToken(foundChild.Pos(), ...)`.

---

### Step 3 — `getRelevantTokens` returns `nil` context and previous tokens

```go
// internal/ls/completions.go ~line 2750
// getRelevantTokens returns contextToken=nil, previousToken=nil
contextToken, previousToken = getRelevantTokens(position, file)
```

`getRelevantTokens` calls `FindPrecedingToken(file, position)` → `nil`.  
Both `contextToken` and `previousToken` remain `nil`.

---

### Step 4 — `tryGetObjectTypeDeclarationCompletionContainer` returns non-nil despite `contextToken == nil`

`location` (from `GetTouchingPropertyName(file, 40)`) = the `foo` `KindIdentifier` node.

```go
// internal/ls/completions.go ~line 4090
func tryGetObjectTypeDeclarationCompletionContainer(
    file *ast.SourceFile,
    contextToken *ast.Node,   // ← nil
    location *ast.Node,
    position int,
) *ast.Node {
    switch location.Kind {
    case ast.KindIdentifier:
        if isFromObjectTypeDeclaration(location) {
            return ast.FindAncestor(location, ast.IsObjectTypeDeclaration)
            // ↑ returns InterfaceDeclaration — early return, contextToken NEVER ACCESSED
        }
    ...
    }
    ...
    // contextToken is only used AFTER this early-return path
}
```

`isFromObjectTypeDeclaration(location)` = `true` (foo → MethodSignature → InterfaceDeclaration),
so the function returns the `InterfaceDeclaration` **without ever touching `contextToken`**.

---

### Step 5 — Nil pointer dereference at line 1447

Back in `tryGetClassLikeCompletionSymbols`:

```go
// internal/ls/completions.go ~line 1437
tryGetClassLikeCompletionSymbols := func() globalsSearch {
    decl := tryGetObjectTypeDeclarationCompletionContainer(file, contextToken, location, position)
    if decl == nil {
        return globalsSearchContinue
    }

    // decl is the InterfaceDeclaration — non-nil, so we continue
    completionKind = CompletionKindMemberLike
    isNewIdentifierLocation = true

    if contextToken.Kind == ast.KindAsteriskToken {  // LINE 1447 — PANIC: nil pointer dereference
```

`contextToken` is `nil`, so `contextToken.Kind` panics.

---

## Comparison with TypeScript Reference

The TypeScript reference implementation in `_submodules/TypeScript` does **not** have the
"single comment optimization" (`IsJSDocSingleCommentNode/IsJSDocSingleCommentNodeList`).
This is a Strada-specific simplification introduced during the Go port.

In the TypeScript reference, `findPrecedingToken(40)` traverses into the JSDoc node's comment
content (because it visits all children including `KindJSDocText`) and returns a non-null token.
Either way, `contextToken` is **non-null**, so `contextToken.kind` does not throw.

The TypeScript reference `tryGetClassLikeCompletionSymbols` has the same unchecked
dereference (`contextToken.kind === 42`), but it is never reached with `contextToken = undefined`
because `findPrecedingToken` always returns something non-null for this position.

This confirms the bug is in the Go port, not in the reference: the "single comment optimization"
inadvertently makes `FindPrecedingToken` return `nil` in cases where TypeScript's
`findPrecedingToken` returns a real token.

---

## Existing Test

**`TestCompletionsAfterJSDoc`** — `internal/fourslash/tests/gen/completionsAfterJSDoc_test.go`

```go
func TestCompletionsAfterJSDoc(t *testing.T) {
    t.Parallel()
    t.Skip()  // ← currently skipped
    defer testutil.RecoverAndFail(t, "Panic on fourslash test")
    const content = `export interface Foo {
  /** JSDoc */
  /**/foo(): void;
}`
    // expects: keyword completion "readonly"
```

Also listed in `internal/fourslash/_scripts/failingTests.txt`.

The `defer testutil.RecoverAndFail(t, "Panic on fourslash test")` confirms this was known
to panic before being skipped.

---

## Two Candidate Fixes

### Fix A — Defensive nil check in `tryGetClassLikeCompletionSymbols` (quick / targeted)

Add a nil guard before `contextToken.Kind` is accessed at line 1447 (and the other
unchecked uses at lines 1461 and 1471):

```go
// internal/ls/completions.go ~line 1447
if contextToken == nil {
    // No context token available (e.g. cursor is right at start of member after JSDoc).
    // Fall through to keyword-based handling without asterisk/semicolon inspection.
    if ast.IsClassLike(decl) {
        keywordFilters = KeywordCompletionFiltersClassElementKeywords
    } else {
        keywordFilters = KeywordCompletionFiltersInterfaceElementKeywords
    }
} else if contextToken.Kind == ast.KindAsteriskToken {
    keywordFilters = KeywordCompletionFiltersNone
} else if ast.IsClassLike(decl) {
    keywordFilters = KeywordCompletionFiltersClassElementKeywords
} else {
    keywordFilters = KeywordCompletionFiltersInterfaceElementKeywords
}
```

The subsequent `if contextToken.Kind == ast.KindSemicolonToken` (line 1461) and
`if contextToken.Kind == ast.KindIdentifier` (line 1471) also need nil guards.

**Pros:** Minimal, targeted, unlikely to break other paths.  
**Cons:** Does not fix the underlying `FindPrecedingToken` inconsistency; other callers
that assume `contextToken` is non-nil may also need fixes in the future.

---

### Fix B — Fix `FindPrecedingTokenEx` to not return nil when JSDoc single-comment traversal yields nothing

When `find(jsDoc)` returns `nil` (due to the single-comment skip), fall back to the
non-JSDoc preceding token search:

```go
// internal/astnav/tokens.go ~line 374 (inside FindPrecedingTokenEx.find)
if jsDoc != nil {
    if !excludeJSDoc {
        if result := find(jsDoc); result != nil {
            return result
        }
        // Fall through: single-comment JSDoc yielded nothing; find non-JSDoc preceding token.
        return findRightmostValidToken(jsDoc.Pos(), sourceFile, n, position, true /*excludeJSDoc*/)
    }
    return findRightmostValidToken(jsDoc.End(), sourceFile, n, position, excludeJSDoc)
}
```

This would make `FindPrecedingToken(file, 40)` return `KindOpenBraceToken` (`{`)
instead of `nil`, which more closely aligns with the TypeScript reference behavior of
returning a non-nil preceding token.

**Pros:** Fixes the root cause; `contextToken` becomes non-nil and the rest of the completion
pipeline works without further changes.  
**Cons:** Changes navigation behavior; needs careful testing to ensure it doesn't break
other positions where `find(jsDoc)` intentionally returns nil.

---

### Fix C — Make `IsJSDocSingleCommentNodeList` transparent inside `findRightmostValidToken`

The `wrappedVisitNodes` guard in `getNodeVisitor` exists to prevent `forEachChild`-style
traversal from visiting single-comment JSDoc text. But `findRightmostValidToken` specifically
needs to find tokens for position-based navigation, and the optimization is breaking it.

A more surgical fix would be to pass a flag to `getNodeVisitor` (or provide a separate
`getNodeVisitorForNavigation`) that skips the `IsJSDocSingleCommentNodeList` guard:

```go
wrappedVisitNodes = func(n *ast.NodeList, v *ast.NodeVisitor) *ast.NodeList {
    if !includeJSDocSingleCommentText && ast.IsJSDocSingleCommentNodeList(n) {
        return n
    }
    return visitNodes(n, v)
}
```

**Pros:** Precisely targeted at the navigation subsystem without affecting `forEachChild`.  
**Cons:** Requires threading a flag through several layers; needs validation that visiting
single-comment JSDoc text in navigation doesn't cause issues elsewhere.

---

## Summary

| Component | Role in Bug |
|-----------|-------------|
| `IsJSDocSingleCommentNode` optimization | Intentionally hides single-comment JSDoc text from child traversal |
| `getNodeVisitor`'s `wrappedVisitNodes` | Skips JSDoc comment NodeList → `findRightmostValidToken` finds no children |
| `FindPrecedingToken(file, 40)` | Returns `nil` instead of `{` or JSDocText token |
| `getRelevantTokens` | Passes `nil` as `contextToken` |
| `tryGetObjectTypeDeclarationCompletionContainer` | Returns non-nil via early-return path that bypasses nil check |
| `tryGetClassLikeCompletionSymbols` line 1447 | Dereferences `contextToken.Kind` without nil check → **panic** |

The cleanest short-term fix is **Fix A** (nil guard in `tryGetClassLikeCompletionSymbols`).
The correct long-term fix is likely **Fix B** or **Fix C** to restore `FindPrecedingToken`
consistency so it never returns nil where a preceding token exists.
