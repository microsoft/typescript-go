// @target: esnext
// @module: preserve  
// @declaration: true

// Function type parenthesized
type Test1<T> = T extends (() => infer R) ? R : never; 

// Call signature in type literal
type Test2<T> = T extends { (): infer R } ? R : never;