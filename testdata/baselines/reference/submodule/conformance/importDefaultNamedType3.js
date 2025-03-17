//// [tests/cases/conformance/externalModules/typeOnly/importDefaultNamedType3.ts] ////

//// [a.ts]
export class A {}

//// [b.ts]
import type from = require('./a');


//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.A = void 0;
class A {
}
exports.A = A;
