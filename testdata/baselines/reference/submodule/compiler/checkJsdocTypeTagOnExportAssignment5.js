//// [tests/cases/compiler/checkJsdocTypeTagOnExportAssignment5.ts] ////

//// [checkJsdocTypeTagOnExportAssignment5.js]

//// [a.js]
/**
 * @typedef {Object} Foo
 * @property {number} a
 * @property {number} b
 */

/** @type {Foo} */
export default { a: 1, b: 1 };

//// [b.js]
import a from "./a";
a;


//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const a_1 = require("./a");
a_1.default;
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = { a: 1, b: 1 };
//// [checkJsdocTypeTagOnExportAssignment5.js]
