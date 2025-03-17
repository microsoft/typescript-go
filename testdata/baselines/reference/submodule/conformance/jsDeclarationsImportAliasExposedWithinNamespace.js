//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsImportAliasExposedWithinNamespace.ts] ////

//// [file.js]
/**
 * @namespace myTypes
 * @global
 * @type {Object<string,*>}
 */
const myTypes = {
    // SOME PROPS HERE
};

/** @typedef {string|RegExp|Array<string|RegExp>} myTypes.typeA */

/**
 * @typedef myTypes.typeB
 * @property {myTypes.typeA}    prop1 - Prop 1.
 * @property {string}           prop2 - Prop 2.
 */

/** @typedef {myTypes.typeB|Function} myTypes.typeC */

export {myTypes};
//// [file2.js]
import {myTypes} from './file.js';

/**
 * @namespace testFnTypes
 * @global
 * @type {Object<string,*>}
 */
const testFnTypes = {
    // SOME PROPS HERE
};

/** @typedef {boolean|myTypes.typeC} testFnTypes.input */

/**
 * @function testFn
 * @description A test function.
 * @param {testFnTypes.input} input - Input.
 * @returns {number|null} Result.
 */
function testFn(input) {
    if (typeof input === 'number') {
        return 2 * input;
    } else {
        return null;
    }
}

export {testFn, testFnTypes};

//// [file2.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.testFnTypes = void 0;
exports.testFn = testFn;
const file_js_1 = require("./file.js");
const testFnTypes = {};
exports.testFnTypes = testFnTypes;
function testFn(input) {
    if (typeof input === 'number') {
        return 2 * input;
    }
    else {
        return null;
    }
}
//// [file.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.myTypes = void 0;
const myTypes = {};
exports.myTypes = myTypes;
