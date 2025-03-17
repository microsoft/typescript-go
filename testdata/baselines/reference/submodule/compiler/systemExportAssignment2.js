//// [tests/cases/compiler/systemExportAssignment2.ts] ////

//// [a.ts]
var a = 10;
export = a;  // Error: export = not allowed in ES6

//// [b.ts]
import * as a from "a";


//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [a.js]
"use strict";
var a = 10;
module.exports = a;
