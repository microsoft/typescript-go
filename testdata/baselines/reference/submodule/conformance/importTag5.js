//// [tests/cases/conformance/jsdoc/importTag5.ts] ////

//// [types.ts]
export interface Foo {
    a: number;
}

//// [foo.js]
/**
 * @import { Foo } from "./types"
 */

/**
 * @param { Foo } foo
 */
function f(foo) {}


//// [types.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [foo.js]
function f(foo) { }
