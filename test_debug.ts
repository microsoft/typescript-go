// Test to understand the flow
type Test = { new(): infer R };

export function testDirectly(): Test {
    return null as any;
}