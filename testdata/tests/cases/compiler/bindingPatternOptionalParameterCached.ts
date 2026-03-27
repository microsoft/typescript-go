// @strict: true
// @noEmit: true

// Usage before declaration triggers getTypeOfSymbol on the parameter,
// caching resolvedType with optionality before the binding pattern check.
const mock: I = {
    m: (_) => {},
};

interface I {
    m({ x }?: { x?: boolean }): void
}
