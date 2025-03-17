//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsTypeReassignmentFromDeclaration.ts] ////

//// [some-mod.d.ts]
interface Item {
    x: string;
}
declare const items: Item[];
export = items;
//// [index.js]
/** @type {typeof import("/some-mod")} */
const items = [];
module.exports = items;

//// [index.js]
const items = [];
module.exports = items;
