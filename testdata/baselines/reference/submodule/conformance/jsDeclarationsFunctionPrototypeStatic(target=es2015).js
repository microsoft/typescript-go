//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsFunctionPrototypeStatic.ts] ////

//// [source.js]
module.exports = MyClass;

function MyClass() {}
MyClass.staticMethod = function() {}
MyClass.prototype.method = function() {}
MyClass.staticProperty = 123;

/**
 * Callback to be invoked when test execution is complete.
 *
 * @callback DoneCB
 * @param {number} failures - Number of failures that occurred.
 */

//// [source.js]
"use strict";
module.exports = MyClass;
function MyClass() { }
MyClass.staticMethod = function () { };
MyClass.prototype.method = function () { };
MyClass.staticProperty = 123;
/**
 * Callback to be invoked when test execution is complete.
 *
 * @callback DoneCB
 * @param {number} failures - Number of failures that occurred.
 */ 


//// [source.d.ts]
export = MyClass;
function MyClass(): void;
declare namespace MyClass {
    var staticMethod: () => void;
}
declare namespace MyClass {
    var staticProperty: number;
}
export type DoneCB = (failures: number) => any;
/**
 * Callback to be invoked when test execution is complete.
 *
 * @callback DoneCB
 * @param {number} failures - Number of failures that occurred.
 */ 


//// [DtsFileErrors]


out/source.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/source.d.ts (1 errors) ====
    export = MyClass;
    function MyClass(): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    declare namespace MyClass {
        var staticMethod: () => void;
    }
    declare namespace MyClass {
        var staticProperty: number;
    }
    export type DoneCB = (failures: number) => any;
    /**
     * Callback to be invoked when test execution is complete.
     *
     * @callback DoneCB
     * @param {number} failures - Number of failures that occurred.
     */ 
    