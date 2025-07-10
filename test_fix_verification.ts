// Test cases that should work with the fix
type Test1 = string extends (x: infer U) => any ? U : never;
type Test2 = string extends { (): infer R } ? R : never;
type Test3 = string extends { new(): infer R } ? R : never;

export function test1(): Test1 { return null as any; }
export function test2(): Test2 { return null as any; }
export function test3(): Test3 { return null as any; }