//// [tests/cases/compiler/expandoFunctionSymbolPropertyJs.ts] ////

//// [types.ts]
export const symb = Symbol();

export interface TestSymb {
  (): void;
  readonly [symb]: boolean;
}

//// [a.js]
import { symb } from "./types";

/**
 * @returns {import("./types").TestSymb}
 */
export function test() {
  function inner() {}
  inner[symb] = true;
  return inner;
}

//// [types.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.symb = void 0;
exports.symb = Symbol();
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.test = test;
const types_1 = require("./types");
function test() {
    function inner() { }
    inner[types_1.symb] = true;
    return inner;
}
