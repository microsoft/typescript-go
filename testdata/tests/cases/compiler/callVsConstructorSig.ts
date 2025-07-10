// @target: esnext
// @module: preserve
// @declaration: true

// Direct comparison: call signature vs constructor signature
type CallSig<T> = T extends { (): infer R } ? R : never;
type ConstructorSig<T> = T extends { new(): infer R } ? R : never;