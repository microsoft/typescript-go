//// [tests/cases/compiler/exportDestructuringIterator.ts] ////

//// [exportDestructuringIterator.ts]
declare function foo(): any;
export const [A, V] = foo();
export const { x, y } = foo();
export const [a = 1, b = 2] = foo();
export const [c, ...d] = foo();
export const [, e, , f] = foo();
export const [[g, h], { i, j: k }] = foo();
export const { m: [n, o], p: { q } } = foo();


//// [exportDestructuringIterator.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.q = exports.o = exports.n = exports.k = exports.i = exports.h = exports.g = exports.f = exports.e = exports.d = exports.c = exports.b = exports.a = exports.y = exports.x = exports.V = exports.A = void 0;
const [A, V] = foo();
exports.A = A, exports.V = V;
const { x, y } = foo();
exports.x = x, exports.y = y;
const [a = 1, b = 2] = foo();
exports.a = a, exports.b = b;
const [c, ...d] = foo();
exports.c = c, exports.d = d;
const [, e, , f] = foo();
exports.e = e, exports.f = f;
const [[g, h], { i, j: k }] = foo();
exports.g = g, exports.h = h, exports.i = i, exports.k = k;
const { m: [n, o], p: { q } } = foo();
exports.n = n, exports.o = o, exports.q = q;
