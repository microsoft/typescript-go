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
Object.defineProperty(exports, "__esModule", { value: true });
export = MyClass;
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
export type DoneCB = (failures: number) ;
/**
 * Callback to be invoked when test execution is complete.
 *
 * @callback DoneCB
 * @param {number} failures - Number of failures that occurred.
 */ 


//// [DtsFileErrors]


out/source.d.ts(1,1): error TS2309: An export assignment cannot be used in a module with other exported elements.
out/source.d.ts(1,10): error TS2304: Cannot find name 'MyClass'.
out/source.d.ts(2,41): error TS1005: '=>' expected.


==== out/source.d.ts (3 errors) ====
    export = MyClass;
    ~~~~~~~~~~~~~~~~~
!!! error TS2309: An export assignment cannot be used in a module with other exported elements.
             ~~~~~~~
!!! error TS2304: Cannot find name 'MyClass'.
    export type DoneCB = (failures: number) ;
                                            ~
!!! error TS1005: '=>' expected.
    /**
     * Callback to be invoked when test execution is complete.
     *
     * @callback DoneCB
     * @param {number} failures - Number of failures that occurred.
     */ 
    