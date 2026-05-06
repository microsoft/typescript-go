//// [tests/cases/conformance/statements/VariableStatements/usingDeclarations/usingDeclarationsDeclarationEmit.1.ts] ////

//// [usingDeclarationsDeclarationEmit.1.ts]
using r1 = { [Symbol.dispose]() {} };
export { r1 };

await using r2 = { async [Symbol.asyncDispose]() {} };
export { r2 };


//// [usingDeclarationsDeclarationEmit.1.js]
using r1 = { [Symbol.dispose]() { } };
export { r1 };
await using r2 = { async [Symbol.asyncDispose]() { } };
export { r2 };


//// [usingDeclarationsDeclarationEmit.1.d.ts]
const r1: {
    [Symbol.dispose](): void;
};
export { r1 };
const r2: {
    [Symbol.asyncDispose](): Promise<void>;
};
export { r2 };


//// [DtsFileErrors]


usingDeclarationsDeclarationEmit.1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== usingDeclarationsDeclarationEmit.1.d.ts (1 errors) ====
    const r1: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        [Symbol.dispose](): void;
    };
    export { r1 };
    const r2: {
        [Symbol.asyncDispose](): Promise<void>;
    };
    export { r2 };
    