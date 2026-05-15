// @checkJs: true
// @allowJs: true
// @strict: true
// @outDir: ./out

// @filename: a.js
/** @typedef {number} NS.T */
/** @typedef {string} NS.U */

/** @type {NS.T} */
const x = 1;

/** @type {NS.U} */
const y = "hello";

// @filename: b.js
/** @typedef {{age: number}} A.B.MyType */

/** @type {A.B.MyType} */
const z = { age: 42 };

// @filename: c.js
/** @callback NS.MyCallback
 * @param {string} name
 * @returns {void}
 */

/** @type {NS.MyCallback} */
const f = (name) => {};
