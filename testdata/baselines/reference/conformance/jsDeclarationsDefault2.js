//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsDefault2.ts] ////

//// [index1.js]
export const _default = class {};

export default 12;
/**
 * @typedef {string | number} default
 */


//// [index1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports._default = void 0;
const _default = class {
};
exports._default = _default;
exports.default = 12;


//// [index1.d.ts]
export declare const _default: {
    new (): {};
};
declare const _default_1: number;
export default _default_1;
export type _default_1 = string | number;
