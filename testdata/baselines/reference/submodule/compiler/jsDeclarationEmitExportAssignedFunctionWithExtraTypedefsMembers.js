//// [tests/cases/compiler/jsDeclarationEmitExportAssignedFunctionWithExtraTypedefsMembers.ts] ////

//// [index.js]
/**
 * @typedef Options
 * @property {string} opt
 */

/**
 * @param {Options} options
 */
module.exports = function loader(options) {}


//// [index.js]
"use strict";
/**
 * @typedef Options
 * @property {string} opt
 */
/**
 * @param {Options} options
 */
module.exports = function loader(options) { };


//// [index.d.ts]
/**
 * @typedef Options
 * @property {string} opt
 */
export type Options = {
    opt: string;
};
/**
 * @param {Options} options
 */
const _default: (options: Options) => void;
export = _default;


//// [DtsFileErrors]


out/index.d.ts(11,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/index.d.ts (1 errors) ====
    /**
     * @typedef Options
     * @property {string} opt
     */
    export type Options = {
        opt: string;
    };
    /**
     * @param {Options} options
     */
    const _default: (options: Options) => void;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export = _default;
    