// @target: esnext
// @module: preserve
// @declaration: true

// This is the exact case from issue #1379
type ExtractReturn<T> = T extends { new(): infer R } ? R : never;