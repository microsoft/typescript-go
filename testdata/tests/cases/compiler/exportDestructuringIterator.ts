// @target: esnext
// @module: commonjs
// @strict: true

declare function foo(): any;
export const [A, V] = foo();
export const { x, y } = foo();
export const [a = 1, b = 2] = foo();
export const [c, ...d] = foo();
