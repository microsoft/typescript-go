//// [tests/cases/compiler/declFileRestParametersOfFunctionAndFunctionType.ts] ////

//// [declFileRestParametersOfFunctionAndFunctionType.ts]
function f1(...args) { }
function f2(x: (...args) => void) { }
function f3(x: { (...args): void }) { }
function f4<T extends (...args) => void>() { }
function f5<T extends { (...args): void }>() { }
var f6 = () => { return [<any>10]; }




//// [declFileRestParametersOfFunctionAndFunctionType.js]
"use strict";
function f1(...args) { }
function f2(x) { }
function f3(x) { }
function f4() { }
function f5() { }
var f6 = () => { return [10]; };


//// [declFileRestParametersOfFunctionAndFunctionType.d.ts]
function f1(...args: any[]): void;
function f2(x: (...args: any[]) => void): void;
function f3(x: {
    (...args: any[]): void;
}): void;
function f4<T extends (...args: any[]) => void>(): void;
function f5<T extends {
    (...args: any[]): void;
}>(): void;
var f6: () => any[];


//// [DtsFileErrors]


declFileRestParametersOfFunctionAndFunctionType.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileRestParametersOfFunctionAndFunctionType.d.ts (1 errors) ====
    function f1(...args: any[]): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function f2(x: (...args: any[]) => void): void;
    function f3(x: {
        (...args: any[]): void;
    }): void;
    function f4<T extends (...args: any[]) => void>(): void;
    function f5<T extends {
        (...args: any[]): void;
    }>(): void;
    var f6: () => any[];
    