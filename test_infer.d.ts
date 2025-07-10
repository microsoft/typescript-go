// Simple test case to understand the infer issue
type ExtractReturn<T> = T extends {
    new ();
} ? R : never;
// Test the problematic case
type Test = ExtractReturn<{
    new ();
}>;
