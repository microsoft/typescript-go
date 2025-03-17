//// [tests/cases/compiler/outModuleConcatUnspecifiedModuleKind.ts] ////

//// [a.ts]
export class A { } // module

//// [b.ts]
var x = 0; // global

//// [b.js]
var x = 0;
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.A = void 0;
class A {
}
exports.A = A;
