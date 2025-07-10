// @target: esnext
// @module: preserve  
// @declaration: true

// Constructor type (like in TypeScript tests)
type ExtractInstanceType<T extends new (...args: any[]) => any> = 
    T extends new (...args: any[]) => infer U ? U : never;

// Constructor signature in type literal (my original test)
type ExtractReturnType<T> = T extends { new(): infer R } ? R : never;