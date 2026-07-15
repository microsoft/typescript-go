//// [tests/cases/compiler/arrowFunctionVisibilityPanic.ts] ////

//// [arrowFunctionVisibilityPanic.ts]
// Regression test for #4629
// This forces the compiler to trigger a symbol-accessibility diagnostic 
// for an ArrowFunction by using a locally scoped (unnameable) type.

export const funcWithPrivateParam = (() => {
    interface PrivateType {
        secret: string;
    }
    return (p: PrivateType) => p;
})();

//// [arrowFunctionVisibilityPanic.js]
// Regression test for #4629
// This forces the compiler to trigger a symbol-accessibility diagnostic 
// for an ArrowFunction by using a locally scoped (unnameable) type.
export const funcWithPrivateParam = (() => {
    return (p) => p;
})();
