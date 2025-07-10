// Test different scenarios to understand the issue
type Test1 = { new(): infer R };
type Test2 = { (): infer R };
type Test3 = (infer R)[];
type Test4 = (x: infer R) => any;

export function test1(): Test1 { return null as any; }
export function test2(): Test2 { return null as any; }
export function test3(): Test3 { return null as any; }
export function test4(): Test4 { return null as any; }