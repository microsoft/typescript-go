# Investigation: Crash with infer + rest/variadic tuple elements

## Bug Summary

A nil pointer dereference (panic) occurs in `internal/checker/inference.go` when type-checking
TypeScript code that uses `infer` with tuple rest (`...T[]`) and variadic (`...infer B extends [any, any]`)
elements in a conditional type, where the source tuple has exactly as many elements as the
variadic element's constraint.

## Reproduction

```ts
type SubTup<T> = T extends [
    ...(infer C)[],
    ...infer B extends [any, any]
] ? B : never;
type Trigger = SubTup<[1, 2]>;
```

This also crashes in the TypeScript compiler (tsc 5.8) with the same root cause:
`TypeError: Cannot read properties of undefined (reading 'aliasSymbol')`.

## Root Cause

The crash is in `inferFromObjectTypes` → `inferFromTypes` in `internal/checker/inference.go`.

When inferring from a source tuple type against a target tuple with two middle elements — a
rest element and a variadic element (or vice versa) — the code calls
`getElementTypeOfSliceOfTupleType()` to extract the rest element's type from the source tuple.
However, this function can return `nil` when the source tuple has no remaining elements for the
rest portion (because the variadic element's `impliedArity` consumes all source elements).

### Detailed walkthrough with `SubTup<[1, 2]>`

1. Target pattern: `[...(infer C)[], ...infer B extends [any, any]]`
   - Element 0: rest element `...(infer C)[]` (flags: `ElementFlagsRest`)
   - Element 1: variadic element `...infer B extends [any, any]` (flags: `ElementFlagsVariadic`)

2. Source: `[1, 2]` — a 2-element tuple

3. This matches the branch at line ~715: "Middle of target is `[...rest, ...T]`"

4. The constraint on B is `[any, any]`, a fixed-size tuple of length 2, so `impliedArity = 2`

5. The code calls:
   ```go
   getElementTypeOfSliceOfTupleType(source, startLength=0, endLength+impliedArity=0+2, false, false)
   ```

6. Inside that function: `length = arity(2) - endSkipCount(2) = 0`.
   Since `index(0) >= length(0)`, it returns `nil`.

7. This `nil` is passed directly as the `source` argument to `inferFromTypes()`.

8. At line 79, `source.alias` dereferences the nil pointer → **panic**.

### Affected code paths (two locations)

1. **Line 710** — `[...T, ...rest]` pattern:
   ```go
   c.inferFromTypes(n, c.getElementTypeOfSliceOfTupleType(source, startLength+impliedArity, endLength, false, false), elementTypes[startLength+1])
   ```

2. **Line 723** — `[...rest, ...T]` pattern:
   ```go
   c.inferFromTypes(n, c.getElementTypeOfSliceOfTupleType(source, startLength, endLength+impliedArity, false, false), elementTypes[startLength])
   ```

Both pass the return value of `getElementTypeOfSliceOfTupleType` directly to `inferFromTypes`
without checking for `nil`. Compare with the correct pattern at line ~736:
```go
restType := c.getElementTypeOfSliceOfTupleType(source, startLength, endLength, false, false)
if restType != nil {
    c.inferFromTypes(n, restType, elementTypes[startLength])
}
```

### TypeScript source equivalent

The TypeScript source at `checker.ts` uses `!` (non-null assertion) at these locations:
```ts
inferFromTypes(getElementTypeOfSliceOfTupleType(source, startLength, endLength + impliedArity)!, elementTypes[startLength]);
```
The `!` suppresses the TypeScript compiler's null check warning but provides no runtime protection.
TypeScript crashes with `TypeError: Cannot read properties of undefined (reading 'aliasSymbol')`.

## Fix

Added nil guards before both `inferFromTypes` calls, consistent with the existing pattern used
in the `middleLength == 1 && ElementFlagsRest` branch at line ~736. When
`getElementTypeOfSliceOfTupleType` returns nil (meaning there are no remaining source elements
for the rest portion), we simply skip the inference for that element — which is semantically
correct since there's nothing to infer from.

## Files Changed

- `internal/checker/inference.go` — Added nil checks at lines 710 and 725
- `testdata/tests/cases/compiler/inferFromTupleRestAndVariadic.ts` — Test case
- `testdata/baselines/reference/compiler/inferFromTupleRestAndVariadic.*` — Baselines
