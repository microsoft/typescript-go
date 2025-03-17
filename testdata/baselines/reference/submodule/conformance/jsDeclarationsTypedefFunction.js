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
let id = 0;
const send = handlers => new Promise((resolve, reject) => {
    handlers[++id] = [resolve, reject];
});
