// Test for the original issue #1379
type ExtractReturn<T> = T extends {
    new ();
} ? R : never;
// Test case that should work
export declare function test(): ExtractReturn<{
    new ();
}>;
export {};
