// Test different infer scenarios
type ExtractReturn<T> = T extends { new(): infer R } ? R : never;
type ExtractParam<T> = T extends (x: infer U) => any ? U : never;
type ExtractArray<T> = T extends (infer U)[] ? U : never;

// More complex cases  
type ExtractCallReturn<T> = T extends { (): infer R } ? R : never;
type ExtractConstructorReturn<T> = T extends new (...args: any[]) => infer U ? U : never;

// Use the types to force them to be emitted
export function test1(): ExtractReturn<{ new(): string }> { return "" as any; }
export function test2(): ExtractParam<(x: number) => void> { return 42 as any; }
export function test3(): ExtractArray<string[]> { return "" as any; }
export function test4(): ExtractCallReturn<{ (): boolean }> { return true as any; }
export function test5(): ExtractConstructorReturn<new () => Date> { return new Date() as any; }