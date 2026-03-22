# Investigation: Issue #3192 — Incorrect `readonly` on tuple in declaration emit

## Problem

After PR #2459 ("Port --isolatedDeclarations related node builder logic"), declaration emit
incorrectly adds a `readonly` modifier to tuples that originate from `as const satisfies` expressions
when the `satisfies` target type has a **mutable** array contextual type.

### Reproduction

```ts
export const obj = {
  array: [
    { n: 1 },
    { n: 2 },
  ],
} as const satisfies { array?: Readonly<{ n: unknown; }>[] }
```

**Expected (matching TypeScript 5.9/6.0):**
```ts
export declare const obj: {
    readonly array: [{
        readonly n: 1;
    }, {
        readonly n: 2;
    }];
};
```

**Actual (tsgo 7.0.0-dev.20260321.1):**
```ts
export declare const obj: {
    readonly array: readonly [{    // ← extra "readonly" on the tuple
        readonly n: 1;
    }, {
        readonly n: 2;
    }];
};
```

## Root Cause

PR #2459 introduced a `PseudoChecker` and `PseudoTypeNodeBuilder` used to try reusing existing
type nodes from source code during declaration emit. The system works as follows:

1. `serializeTypeForDeclaration` (in `nodebuilderimpl.go:2073`) tries to use the pseudochecker
   to get the type of a declaration before falling back to the full type serializer.

2. `PseudoTypeTuple` (in `pseudochecker/type.go:195`) represents tuples from `as const` array
   literals but has **no `readonly` field** — it implicitly assumes all such tuples are readonly.

3. `pseudoTypeToNode` (in `pseudotypenodebuilder.go:129-139`) **always** wraps tuple type nodes
   with `NewTypeOperatorNode(ast.KindReadonlyKeyword, ...)`, with the comment "pseudo-tuples are
   implicitly `readonly` since they originate from `as const` contexts".

4. `pseudoTypeEquivalentToType` (in `pseudotypenodebuilder.go:420-440`) checks whether a
   pseudo-tuple matches a real checker tuple by comparing element counts and element types, but
   **does not check the `readonly` flag** on the real tuple.

The bug manifests when:
- `as const` creates the array in a const context
- `satisfies T` provides a contextual type where T has a **mutable** array (e.g., `Readonly<X>[]`)
- The checker's `checkArrayLiteral` correctly determines the tuple is NOT readonly (because the
  contextual type is mutable-array-like: `isMutableArrayLikeType` returns `true`)
- But the pseudochecker assumes it IS readonly, and the equivalence check doesn't catch the
  mismatch
- The incorrect pseudo-type node (with readonly) is used instead of falling back to the correct
  type serialization path

## Fix

Added a `readonly` check in `pseudoTypeEquivalentToType` for the `PseudoTypeKindTuple` case.
Since pseudo-tuples always emit as readonly, they can only correctly match real tuples that are
also readonly. When the real tuple is non-readonly, the equivalence check now returns `false`,
causing the code to fall back to the regular `typeToTypeNode` path which correctly handles the
readonly flag based on the actual type.

**Changed file:** `internal/checker/pseudotypenodebuilder.go`

```go
// Before (line 420-430):
case pseudochecker.PseudoTypeKindTuple:
    // ... existing checks for element flags ...

// After:
case pseudochecker.PseudoTypeKindTuple:
    // ... check tupleTarget.readonly ...
    if !tupleTarget.readonly {
        return false  // Pseudo-tuples always emit as readonly
    }
    // ... existing checks for element flags ...
```

## Related

- PR #2459: Port --isolatedDeclarations related node builder logic
- TypeScript PR #55229 + #55522: Correct readonly tuple behavior referenced by @Andarist
- Comment #4105762522 for the original test case
