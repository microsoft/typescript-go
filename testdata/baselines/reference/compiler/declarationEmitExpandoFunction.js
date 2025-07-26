//// [tests/cases/compiler/declarationEmitExpandoFunction.ts] ////

//// [declarationEmitExpandoFunction.ts]
export function A() {
    return 'A';
}

export enum B {
    C
}

A.B = B;


//// [declarationEmitExpandoFunction.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.B = void 0;
exports.A = A;
function A() {
    return 'A';
}
var B;
(function (B) {
    B[B["C"] = 0] = "C";
})(B || (exports.B = B = {}));
A.B = B;


//// [declarationEmitExpandoFunction.d.ts]
export declare function A(): string;
export declare namespace A {
    var B: typeof import("./declarationEmitExpandoFunction").B;
}
export declare enum B {
    C = 0
}
