//// [tests/cases/compiler/filesEmittingIntoSameOutputWithOutOption.ts] ////

//// [a.ts]
export class c {
}

//// [b.ts]
function foo() {
}


//// [b.js]
function foo() {
}
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.c = void 0;
class c {
}
exports.c = c;
