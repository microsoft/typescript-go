//// [tests/cases/compiler/es6ExportAssignment2.ts] ////

//// [a.ts]
var a = 10;
export = a;  // Error: export = not allowed in ES6

//// [b.ts]
import * as a from "a";


//// [b.js]
export {};
//// [a.js]
var a = 10;
export {};
