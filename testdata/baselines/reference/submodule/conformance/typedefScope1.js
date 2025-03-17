//// [tests/cases/conformance/jsdoc/typedefScope1.ts] ////

//// [typedefScope1.js]
function B1() {
    /** @typedef {number} B */
    /** @type {B} */
    var ok1 = 0;
}

function B2() {
    /** @typedef {string} B */
    /** @type {B} */
    var ok2 = 'hi';
}

/** @type {B} */
var notOK = 0;


//// [typedefScope1.js]
function B1() {
    var ok1 = 0;
}
function B2() {
    var ok2 = 'hi';
}
var notOK = 0;
