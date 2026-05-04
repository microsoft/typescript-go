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
[exports.A, exports.V] = foo();
({ x: exports.x, y: exports.y } = foo());
[exports.a = 1, exports.b = 2] = foo();
[exports.c, ...exports.d] = foo();
