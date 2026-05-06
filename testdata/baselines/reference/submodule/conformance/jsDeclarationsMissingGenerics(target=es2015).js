//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsMissingGenerics.ts] ////

//// [file.js]
/**
 * @param {Array} x
 */
function x(x) {}
/**
 * @param {Promise} x
 */
function y(x) {}

//// [file.js]
"use strict";
/**
 * @param {Array} x
 */
function x(x) { }
/**
 * @param {Promise} x
 */
function y(x) { }


//// [file.d.ts]
/**
 * @param {Array} x
 */
function x(x: Array): void;
/**
 * @param {Promise} x
 */
function y(x: Promise): void;


//// [DtsFileErrors]


out/file.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
out/file.d.ts(4,15): error TS2314: Generic type 'Array<T>' requires 1 type argument(s).
out/file.d.ts(8,15): error TS2314: Generic type 'Promise<T>' requires 1 type argument(s).


==== out/file.d.ts (3 errors) ====
    /**
     * @param {Array} x
     */
    function x(x: Array): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
                  ~~~~~
!!! error TS2314: Generic type 'Array<T>' requires 1 type argument(s).
    /**
     * @param {Promise} x
     */
    function y(x: Promise): void;
                  ~~~~~~~
!!! error TS2314: Generic type 'Promise<T>' requires 1 type argument(s).
    