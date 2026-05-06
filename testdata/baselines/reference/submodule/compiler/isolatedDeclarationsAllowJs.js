//// [tests/cases/compiler/isolatedDeclarationsAllowJs.ts] ////

//// [file1.ts]
export var x;
//// [file2.js]
export var y;

//// [file1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.x = void 0;


//// [file1.d.ts]
export var x: any;
//// [file2.d.ts]
export var y: any;
