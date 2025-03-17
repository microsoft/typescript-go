//// [tests/cases/compiler/isolatedDeclarationsAllowJs.ts] ////

//// [file1.ts]
export var x;
//// [file2.js]
export var y;

//// [file2.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.y = void 0;
//// [file1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.x = void 0;
