type Test1 = string extends (x: infer U) => any ? U : never;
type Test2 = string extends {
    (): infer R;
} ? R : never;
type Test3 = string extends {
    new (): infer R;
} ? R : never;
export declare function test1(): Test1;
export declare function test2(): Test2;
export declare function test3(): Test3;
export {};
