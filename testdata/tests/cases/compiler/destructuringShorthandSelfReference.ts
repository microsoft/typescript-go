// @strict: true
// @noEmit: true

// Test for stack overflow fix when a shorthand property in an object literal
// references a variable being declared in the same destructuring pattern.
// See: https://github.com/microsoft/TypeScript/issues/62993

const { c, f }: string | number | symbol = { c: 0, f };
