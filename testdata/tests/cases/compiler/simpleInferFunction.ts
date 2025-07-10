// @target: esnext
// @module: preserve
// @declaration: true

// Simple function parameter case
type ExtractParam<T> = T extends (x: infer R) => any ? R : never;