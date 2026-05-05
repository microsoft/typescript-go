//// [tests/cases/compiler/exportDestructuringIterator.ts] ////

//// [exportDestructuringIterator.ts]
declare function foo(): any;
export const [A, V] = foo();
export const { x, y } = foo();
export const [a = 1, b = 2] = foo();
export const [c, ...d] = foo();


//// [exportDestructuringIterator.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.d = exports.c = exports.b = exports.a = exports.y = exports.x = exports.V = exports.A = void 0;
const [A, V] = foo();
exports.A = A, exports.V = V;
const { x, y } = foo();
exports.x = x, exports.y = y;
const [a = 1, b = 2] = foo();
exports.a = a, exports.b = b;
const [c, ...d] = foo();
exports.c = c, exports.d = d;
