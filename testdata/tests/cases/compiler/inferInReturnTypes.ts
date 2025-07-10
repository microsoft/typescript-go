// @target: esnext
// @module: preserve  
// @declaration: true

// Test infer in function return type
type ExtractFunctionReturn<T> = T extends () => infer R ? R : never;

// Test infer in constructor return type  
type ExtractConstructorReturn<T> = T extends { new(): infer R } ? R : never;