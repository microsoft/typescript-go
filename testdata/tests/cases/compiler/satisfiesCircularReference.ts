// @strict: true

// Circular reference after satisfies should produce errors, not panic
const f = () => 42 satisfies typeof f

// These cases should work without errors
const g = () => g satisfies typeof g
