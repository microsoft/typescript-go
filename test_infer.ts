// Simple test case to understand the infer issue
type ExtractReturn<T> = T extends { new(): infer R } ? R : never;

// Test the problematic case
type Test = ExtractReturn<{ new(): string }>;