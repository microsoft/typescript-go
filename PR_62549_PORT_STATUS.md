# Port Status: microsoft/TypeScript#62549

## Summary
**Status: ✅ ALREADY PORTED**

PR #62549 from microsoft/TypeScript has already been successfully ported to the Go codebase.

## Details

### Original Change
The PR "Consistently resolve to the `errorType` on `arguments` with error" changed the TypeScript checker to pass `ignoreArrowFunctions: true` when checking if `arguments` is used in a property initializer or class static block.

**TypeScript change (src/compiler/checker.ts):**
```diff
- if (isInPropertyInitializerOrClassStaticBlock(node)) {
+ if (isInPropertyInitializerOrClassStaticBlock(node, /*ignoreArrowFunctions*/ true)) {
```

### Go Implementation
The equivalent change is already present in the Go codebase:

**File:** `internal/checker/checker.go`, line 10599
```go
if c.isInPropertyInitializerOrClassStaticBlock(node, true /*ignoreArrowFunctions*/) {
    c.error(node, diagnostics.X_arguments_cannot_be_referenced_in_property_initializers_or_class_static_initialization_blocks)
    return c.errorType
}
```

### Test Results
- ✅ All tests pass
- ✅ The specific test `argumentsUsedInClassFieldInitializerOrStaticInitializationBlock` passes
- ✅ The testdata baseline shows the correct behavior (returning `any` type instead of `IArguments` for `arguments` in arrow functions within property initializers)

### Baseline Differences
The testdata contains `.diff` files showing differences between the Go implementation and the TypeScript submodule baselines. This is expected because:
1. The Go implementation has the change
2. The TypeScript submodule baselines haven't been regenerated yet (shallow clone)
3. The test framework correctly handles these differences

## History
The change from PR #62549 was actually ported **before** it was merged into TypeScript!

- **Oct 5, 2025**: @Andarist opened PR #1828 to port microsoft/TypeScript#48172 which added the `ignoreArrowFunctions` parameter
- **In PR #1828**: The Go implementation was already passing `true` for `ignoreArrowFunctions` in the arguments check
- **Oct 5, 2025**: @Andarist opened microsoft/TypeScript#62549 to make TypeScript match the Go behavior
- **Oct 5, 2025**: PR #62549 was merged into TypeScript

As noted in the [PR #1828 review](https://github.com/microsoft/typescript-go/pull/1828#discussion_r2404688959):
> "this has some seemingly unwanted diffs but I think this is basically desired and I opened a PR to Strada to close this gap: microsoft/TypeScript#62549"

## Conclusion
No action needed. The change has been correctly ported and is functioning as expected. In fact, the Go port led the way, and TypeScript followed!
