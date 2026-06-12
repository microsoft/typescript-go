//// [tests/cases/compiler/inferenceUnionArrayDepth.ts] ////

//// [inferenceUnionArrayDepth.ts]
// Regression test for https://github.com/microsoft/typescript-go/issues/1789
// and https://github.com/microsoft/typescript-go/issues/3370
// Inference should correctly infer T from T[] | T[][] union parameters.

declare function flat<T>(args: T[] | T[][]): T;

// Case 1: Union type (issue #1789)
type Value = 1 | 2;
declare const n: Value[] | Value[][];
const result1: Value = flat(n); // Should infer T = Value, not T = Value[]

// Case 2: Object type (issue #3370)
type TG = { a: string };

function isNestedArray<T>(arr: T[] | T[][]): arr is T[][] {
    return Array.isArray(arr) && Array.isArray(arr[0]);
}

function convert(controls: TG[] | TG[][]): TG[][] {
    if (isNestedArray(controls)) {
        return controls;
    } else {
        return [controls];
    }
}

// Case 3: Primitive type (should already work)
declare const s: string[] | string[][];
const result3: string = flat(s);


//// [inferenceUnionArrayDepth.js]
"use strict";
// Regression test for https://github.com/microsoft/typescript-go/issues/1789
// and https://github.com/microsoft/typescript-go/issues/3370
// Inference should correctly infer T from T[] | T[][] union parameters.
const result1 = flat(n); // Should infer T = Value, not T = Value[]
function isNestedArray(arr) {
    return Array.isArray(arr) && Array.isArray(arr[0]);
}
function convert(controls) {
    if (isNestedArray(controls)) {
        return controls;
    }
    else {
        return [controls];
    }
}
const result3 = flat(s);
