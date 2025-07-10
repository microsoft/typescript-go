// @target: esnext
// @module: preserve  
// @declaration: true

// Simple case that should use the simple branch
type ExtractParam<T> = T extends (x: infer U) => any ? U : never;