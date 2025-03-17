//// [tests/cases/conformance/externalModules/typeOnly/preserveValueImports.ts] ////

//// [a.ts]
export default {};
export const b = 0;
export const c = 1;
export interface D {}

//// [b.ts]
import a, { b, c, D } from "./a";

//// [c.ts]
import * as a from "./a";

//// [d.ts]
export = {};

//// [e.ts]
import D = require("./d");
import DD = require("./d");
DD;

//// [f.ts]
import type a from "./a";
import { b, c } from "./a";
b;


//// [f.js]
import { b } from "./a";
b;
//// [e.js]
DD;
export {};
//// [d.js]
export {};
//// [c.js]
export {};
//// [b.js]
export {};
//// [a.js]
export default {};
export const b = 0;
export const c = 1;
