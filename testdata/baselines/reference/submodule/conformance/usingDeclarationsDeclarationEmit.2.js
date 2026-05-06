//// [tests/cases/conformance/statements/VariableStatements/usingDeclarations/usingDeclarationsDeclarationEmit.2.ts] ////

//// [usingDeclarationsDeclarationEmit.2.ts]
using r1 = { [Symbol.dispose]() {} };
export type R1 = typeof r1;

await using r2 = { async [Symbol.asyncDispose]() {} };
export type R2 = typeof r2;


//// [usingDeclarationsDeclarationEmit.2.js]
using r1 = { [Symbol.dispose]() { } };
await using r2 = { async [Symbol.asyncDispose]() { } };
export {};


//// [usingDeclarationsDeclarationEmit.2.d.ts]
const r1: {
    [Symbol.dispose](): void;
};
export type R1 = typeof r1;
const r2: {
    [Symbol.asyncDispose](): Promise<void>;
};
export type R2 = typeof r2;
export {};


//// [DtsFileErrors]


usingDeclarationsDeclarationEmit.2.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== usingDeclarationsDeclarationEmit.2.d.ts (1 errors) ====
    const r1: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        [Symbol.dispose](): void;
    };
    export type R1 = typeof r1;
    const r2: {
        [Symbol.asyncDispose](): Promise<void>;
    };
    export type R2 = typeof r2;
    export {};
    