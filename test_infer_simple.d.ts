// Test different scenarios to understand the issue
type Test1 = {
    new ();
};
type Test2 = {
    (): infer R;
};
type Test3 = (infer R)[];
type Test4 = (x: infer R) => any;
export declare function test1(): Test1;
export declare function test2(): Test2;
export declare function test3(): Test3;
export declare function test4(): Test4;
export {};
