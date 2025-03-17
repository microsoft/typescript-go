//// [tests/cases/compiler/nodeResolution1.ts] ////

//// [a.ts]
export var x = 1;

//// [b.ts]
import y = require("./a");

//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.x = void 0;
exports.x = 1;
