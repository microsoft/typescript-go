//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsTypeAliases.ts] ////

//// [index.js]
export {}; // flag file as module
/**
 * @typedef {string | number | symbol} PropName 
 */

/**
 * Callback
 *
 * @callback NumberToStringCb
 * @param {number} a
 * @returns {string}
 */

/**
 * @template T
 * @typedef {T & {name: string}} MixinName 
 */

/**
 * Identity function
 *
 * @template T
 * @callback Identity
 * @param {T} x
 * @returns {T}
 */

//// [mixed.js]
/**
 * @typedef {{x: string} | number | LocalThing | ExportedThing} SomeType
 */
/**
 * @param {number} x
 * @returns {SomeType}
 */
function doTheThing(x) {
    return {x: ""+x};
}
class ExportedThing {
    z = "ok"
}
module.exports = {
    doTheThing,
    ExportedThing,
};
class LocalThing {
    y = "ok"
}


//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [mixed.js]
function doTheThing(x) {
    return { x: "" + x };
}
class ExportedThing {
    z = "ok";
}
module.exports = {
    doTheThing,
    ExportedThing,
};
class LocalThing {
    y = "ok";
}
