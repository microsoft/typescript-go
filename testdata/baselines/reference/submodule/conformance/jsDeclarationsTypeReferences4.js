//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsTypeReferences4.ts] ////

//// [index.d.ts]
declare module "fs" {
    export class Something {}
}
//// [index.js]
/// <reference types="node" />
export const Something = 2; // to show conflict that can occur
// @ts-ignore
export namespace A {
    // @ts-ignore
    export namespace B {
        const Something = require("fs").Something;
        const thing = new Something();
        // @ts-ignore
        export { thing };
    }
}


//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.A = exports.Something = void 0;
exports.Something = 2;
var A;
(function (A) {
    let B;
    (function (B) {
        const Something = require("fs").Something;
        const thing = new Something();
        export { thing };
    })(B = A.B || (A.B = {}));
})(A || (exports.A = A = {}));
