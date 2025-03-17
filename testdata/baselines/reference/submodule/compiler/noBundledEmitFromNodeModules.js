//// [tests/cases/compiler/noBundledEmitFromNodeModules.ts] ////

//// [index.ts]
export class C {}

//// [a.ts]
import { C } from "projB";


//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.C = void 0;
class C {
}
exports.C = C;
