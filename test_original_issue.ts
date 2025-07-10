// Test for the original issue #1379
type ExtractReturn<T> = T extends { new(): infer R } ? R : never;

// Test case that should work
export function test(): ExtractReturn<{ new(): string }> {
    return "test" as any;
}