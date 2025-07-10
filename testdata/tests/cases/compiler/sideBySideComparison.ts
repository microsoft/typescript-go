// @target: esnext
// @module: preserve  
// @declaration: true

// Call signature in type literal - should work
type ExtractCallReturn<T> = T extends { (): infer R } ? R : never;

// Constructor signature in type literal - should be fixed
type ExtractConstructReturn<T> = T extends { new(): infer R } ? R : never;

// Test both with same usage
declare const callTest: ExtractCallReturn<() => string>;
declare const constructTest: ExtractConstructReturn<{ new(): string }>;