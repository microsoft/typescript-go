// @target: esnext
// @module: preserve  
// @declaration: true

// Simple constructor signature with infer
type SimpleConstructor<T> = T extends { new(): infer R } ? R : never;

// Constructor with parameters
type ConstructorWithParam<T> = T extends { new(x: any): infer R } ? R : never;