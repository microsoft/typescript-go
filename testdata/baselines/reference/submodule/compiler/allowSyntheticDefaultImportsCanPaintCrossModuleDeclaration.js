//// [tests/cases/compiler/allowSyntheticDefaultImportsCanPaintCrossModuleDeclaration.ts] ////

//// [color.ts]
interface Color {
    c: string;
}
export default Color;
//// [file1.ts]
import Color from "./color";
export declare function styled(): Color;
//// [file2.ts]
import { styled }  from "./file1";
export const A = styled();

//// [file2.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.A = void 0;
const file1_1 = require("./file1");
exports.A = (0, file1_1.styled)();
//// [file1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [color.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
