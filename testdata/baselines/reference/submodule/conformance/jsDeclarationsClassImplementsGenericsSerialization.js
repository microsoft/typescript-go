//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsClassImplementsGenericsSerialization.ts] ////

//// [interface.ts]
export interface Encoder<T> {
    encode(value: T): Uint8Array
}
//// [lib.js]
/**
 * @template T
 * @implements {IEncoder<T>}
 */
export class Encoder {
    /**
     * @param {T} value 
     */
    encode(value) {
        return new Uint8Array(0)
    }
}


/**
 * @template T
 * @typedef {import('./interface').Encoder<T>} IEncoder
 */

//// [interface.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [lib.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Encoder = void 0;
class Encoder {
    encode(value) {
        return new Uint8Array(0);
    }
}
exports.Encoder = Encoder;
