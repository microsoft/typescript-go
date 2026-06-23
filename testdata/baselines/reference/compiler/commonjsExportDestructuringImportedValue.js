//// [tests/cases/compiler/commonjsExportDestructuringImportedValue.ts] ////

//// [enum.ts]
export class CodePriceType {
    static A = "a";
    static B = "b";
}

//// [repro.ts]
import { CodePriceType } from "./enum";
export const { A, B } = CodePriceType;


//// [enum.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CodePriceType = void 0;
class CodePriceType {
    static A = "a";
    static B = "b";
}
exports.CodePriceType = CodePriceType;
//// [repro.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.B = exports.A = void 0;
const enum_1 = require("./enum");
({ A: exports.A, B: exports.B } = CodePriceType);
