//// [tests/cases/conformance/externalModules/typeOnly/exportDeclaration_value.ts] ////

//// [a.ts]
const A = {};
export type { A };
export const AA = {};

//// [b.ts]
export type { AA } from './a';


//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.AA = void 0;
const A = {};
exports.AA = {};
