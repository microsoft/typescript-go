// @declaration: true
// Test case to isolate the constructor signature issue
type ExtractConstructorReturn<T> = T extends { new(): infer R } ? R : never;

export function test(): ExtractConstructorReturn<{ new(): string }> {
    return "" as any;
}