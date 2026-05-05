//// [tests/cases/compiler/exportObjectRest.ts] ////

//// [exportObjectRest.ts]
export const { x, ...rest } = { x: 'x', y: 'y' };

//// [exportObjectRest.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.rest = exports.x = void 0;
const { x, ...rest } = { x: 'x', y: 'y' };
exports.x = x, exports.rest = rest;
