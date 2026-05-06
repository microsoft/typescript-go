//// [tests/cases/compiler/declarationEmitDestructuring4.ts] ////

//// [declarationEmitDestructuring4.ts]
// For an array binding pattern with empty elements,
// we will not make any modification and will emit
// the similar binding pattern users' have written
function baz([]) { }
function baz1([] = [1,2,3]) { }
function baz2([[]] = [[1,2,3]]) { }

function baz3({}) { }
function baz4({} = { x: 10 }) { }



//// [declarationEmitDestructuring4.js]
"use strict";
// For an array binding pattern with empty elements,
// we will not make any modification and will emit
// the similar binding pattern users' have written
function baz([]) { }
function baz1([] = [1, 2, 3]) { }
function baz2([[]] = [[1, 2, 3]]) { }
function baz3({}) { }
function baz4({} = { x: 10 }) { }


//// [declarationEmitDestructuring4.d.ts]
function baz([]: Iterable<any, void, undefined>): void;
function baz1([]?: number[]): void;
function baz2([[]]?: [number[]]): void;
function baz3({}: {}): void;
function baz4({}?: {
    x: number;
}): void;


//// [DtsFileErrors]


declarationEmitDestructuring4.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitDestructuring4.d.ts (1 errors) ====
    function baz([]: Iterable<any, void, undefined>): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function baz1([]?: number[]): void;
    function baz2([[]]?: [number[]]): void;
    function baz3({}: {}): void;
    function baz4({}?: {
        x: number;
    }): void;
    