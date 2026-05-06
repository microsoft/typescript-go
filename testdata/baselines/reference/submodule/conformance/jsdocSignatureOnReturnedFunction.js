//// [tests/cases/conformance/jsdoc/jsdocSignatureOnReturnedFunction.ts] ////

//// [jsdocSignatureOnReturnedFunction.js]
function f1() {
    /**
     * @param {number} a
     * @param {number} b
     * @returns {number}
     */
    return (a, b) => {
        return a + b;
    }
}

function f2() {
    /**
     * @param {number} a
     * @param {number} b
     * @returns {number}
     */
    return function (a, b){
        return a + b;
    }
}

function f3() {
    /** @type {(a: number, b: number) => number} */
    return (a, b) => {
        return a + b;
    }
}

function f4() {
    /** @type {(a: number, b: number) => number} */
    return function (a, b){
        return a + b;
    }
}


//// [jsdocSignatureOnReturnedFunction.js]
"use strict";
function f1() {
    /**
     * @param {number} a
     * @param {number} b
     * @returns {number}
     */
    return (a, b) => {
        return a + b;
    };
}
function f2() {
    /**
     * @param {number} a
     * @param {number} b
     * @returns {number}
     */
    return function (a, b) {
        return a + b;
    };
}
function f3() {
    /** @type {(a: number, b: number) => number} */
    return (a, b) => {
        return a + b;
    };
}
function f4() {
    /** @type {(a: number, b: number) => number} */
    return function (a, b) {
        return a + b;
    };
}


//// [jsdocSignatureOnReturnedFunction.d.ts]
function f1(): (a: number, b: number) => number;
function f2(): (a: number, b: number) => number;
function f3(): (a: number, b: number) => number;
function f4(): (a: number, b: number) => number;


//// [DtsFileErrors]


out/jsdocSignatureOnReturnedFunction.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/jsdocSignatureOnReturnedFunction.d.ts (1 errors) ====
    function f1(): (a: number, b: number) => number;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function f2(): (a: number, b: number) => number;
    function f3(): (a: number, b: number) => number;
    function f4(): (a: number, b: number) => number;
    