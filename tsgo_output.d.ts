// Test different infer scenarios
type ExtractReturn<T> = T extends {
    new ();
} ? R : never;
type ExtractParam<T> = T extends (x: infer U) => any ? U : never;
type ExtractArray<T> = T extends (infer U)[] ? U : never;
// More complex cases  
type ExtractCallReturn<T> = T extends {
    (): infer R;
} ? R : never;
type ExtractConstructorReturn<T> = T extends new (...args: any[]) => infer U ? U : never;
// Use the types to force them to be emitted
export declare function test1(): ExtractReturn<{
    new ();
}>;
export declare function test2(): ExtractParam<(x: number) => void>;
export declare function test3(): ExtractArray<string[]>;
export declare function test4(): ExtractCallReturn<{
    (): boolean;
}>;
export declare function test5(): ExtractConstructorReturn<new () => Date>;
export {};
