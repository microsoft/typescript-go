//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsTypedefFunction.ts] ////

//// [foo.js]
/**
 * @typedef {{
 *   [id: string]: [Function, Function];
 * }} ResolveRejectMap
 */

let id = 0

/**
 * @param {ResolveRejectMap} handlers
 * @returns {Promise<any>}
 */
const send = handlers => new Promise((resolve, reject) => {
  handlers[++id] = [resolve, reject]
})

//// [foo.js]
"use strict";
/**
 * @typedef {{
 *   [id: string]: [Function, Function];
 * }} ResolveRejectMap
 */
let id = 0;
/**
 * @param {ResolveRejectMap} handlers
 * @returns {Promise<any>}
 */
const send = handlers => new Promise((resolve, reject) => {
    handlers[++id] = [resolve, reject];
});


//// [foo.d.ts]
/**
 * @typedef {{
 *   [id: string]: [Function, Function];
 * }} ResolveRejectMap
 */
type ResolveRejectMap = {
    [id: string]: [Function, Function];
};
let id: number;
/**
 * @param {ResolveRejectMap} handlers
 * @returns {Promise<any>}
 */
const send: (handlers: ResolveRejectMap) => Promise<any>;


//// [DtsFileErrors]


out/foo.d.ts(9,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/foo.d.ts (1 errors) ====
    /**
     * @typedef {{
     *   [id: string]: [Function, Function];
     * }} ResolveRejectMap
     */
    type ResolveRejectMap = {
        [id: string]: [Function, Function];
    };
    let id: number;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    /**
     * @param {ResolveRejectMap} handlers
     * @returns {Promise<any>}
     */
    const send: (handlers: ResolveRejectMap) => Promise<any>;
    