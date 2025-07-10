type ExtractReturn<T> = T extends {
    new (): infer R;
} ? R : never;
type ExtractParam<T> = T extends (x: infer U) => any ? U : never;
type ExtractArray<T> = T extends (infer U)[] ? U : never;
type ExtractCallReturn<T> = T extends {
    (): infer R;
} ? R : never;
type ExtractConstructorReturn<T> = T extends new (...args: any[]) => infer U ? U : never;
export declare function test1(): ExtractReturn<{
    new (): string;
}>;
export declare function test2(): ExtractParam<(x: number) => void>;
export declare function test3(): ExtractArray<string[]>;
export declare function test4(): ExtractCallReturn<{
    (): boolean;
}>;
export declare function test5(): ExtractConstructorReturn<new () => Date>;
export {};
