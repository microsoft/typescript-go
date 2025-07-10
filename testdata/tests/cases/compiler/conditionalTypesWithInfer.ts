// @target: esnext
// @module: preserve
// @declaration: true

// Basic conditional type with infer
type ExtractReturn<T> = T extends { new(): infer R } ? R : never;

// More complex examples
type ExtractFirstParam<T> = T extends (first: infer U, ...rest: any[]) => any ? U : never;

type UnwrapArray<T> = T extends (infer U)[] ? U : T;

type ExtractPromiseValue<T> = T extends Promise<infer V> ? V : T;

// Multiple infer in same conditional
type ExtractFunction<T> = T extends (arg: infer A) => infer R ? { arg: A; return: R } : never;

// Nested conditionals with infer
type DeepExtract<T> = T extends Promise<infer U> 
  ? U extends Array<infer V> 
    ? V 
    : U 
  : T;

// Usage examples
declare const extractReturn: ExtractReturn<{ new(): string }>;
declare const extractParam: ExtractFirstParam<(x: number, y: string) => void>;
declare const unwrapArray: UnwrapArray<string[]>;
declare const extractPromise: ExtractPromiseValue<Promise<boolean>>;
declare const extractFunction: ExtractFunction<(x: number) => string>;
declare const deepExtract: DeepExtract<Promise<Array<Date>>>;