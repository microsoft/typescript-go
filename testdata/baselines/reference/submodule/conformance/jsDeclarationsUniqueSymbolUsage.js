//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsUniqueSymbolUsage.ts] ////

//// [a.js]
export const kSymbol = Symbol("my-symbol");

/**
 * @typedef {{[kSymbol]: true}} WithSymbol
 */
//// [b.js]
/**
 * @returns {import('./a').WithSymbol} 
 * @param {import('./a').WithSymbol} value 
 */
export function b(value) {
    return value;
}


//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.b = b;
function b(value) {
    return value;
}
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.kSymbol = void 0;
exports.kSymbol = Symbol("my-symbol");
