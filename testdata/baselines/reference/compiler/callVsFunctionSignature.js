//// [tests/cases/compiler/callVsFunctionSignature.ts] ////

//// [callVsFunctionSignature.ts]
// Function type parenthesized
type Test1<T> = T extends (() => infer R) ? R : never; 

// Call signature in type literal
type Test2<T> = T extends { (): infer R } ? R : never;

//// [callVsFunctionSignature.js]


//// [callVsFunctionSignature.d.ts]
// Function type parenthesized
type Test1<T> = T extends (() => infer R) ? R : never;
// Call signature in type literal
type Test2<T> = T extends {
    (): infer R;
} ? R : never;
