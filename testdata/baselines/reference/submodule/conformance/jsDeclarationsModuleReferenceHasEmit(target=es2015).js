//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsModuleReferenceHasEmit.ts] ////

//// [index.js]
/**
 * @module A
 */
class A {}


/**
 * Target element
 * @type {module:A}
 */
export let el = null;

export default A;

//// [index.js]
/**
 * @module A
 */
class A {
}
/**
 * Target element
 * @type {module:A}
 */
export let el = null;
export default A;


//// [index.d.ts]
/**
 * @module A
 */
class A {
}
/**
 * Target element
 * @type {module:A}
 */
export let el: module;
export default A;
