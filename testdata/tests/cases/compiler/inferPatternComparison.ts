// @target: esnext
// @module: preserve
// @declaration: true

// Test various patterns to understand where the issue occurs

// 1. Function parameter (works)
type Test1<T> = T extends (x: infer R) => any ? R : never;

// 2. Function return type (should work) 
type Test2<T> = T extends () => infer R ? R : never;

// 3. Constructor signature return type (broken)
type Test3<T> = T extends { new(): infer R } ? R : never;

// 4. Constructor parameter (might be broken too)
type Test4<T> = T extends { new(x: infer R): any } ? R : never;

// 5. Call signature return type (should work like function)
type Test5<T> = T extends { (): infer R } ? R : never;

// 6. Call signature parameter (should work like function)  
type Test6<T> = T extends { (x: infer R): any } ? R : never;