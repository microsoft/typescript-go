//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsImportNamespacedType.ts] ////

//// [file.js]
import { dummy } from './mod1'
/** @type {import('./mod1').Dotted.Name} - should work */
var dot2

//// [mod1.js]
/** @typedef {number} Dotted.Name */
export var dummy = 1


//// [mod1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.dummy = void 0;
exports.dummy = 1;
//// [file.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const mod1_1 = require("./mod1");
var dot2;
