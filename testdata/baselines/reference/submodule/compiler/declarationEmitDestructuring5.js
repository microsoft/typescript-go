//// [tests/cases/compiler/declarationEmitDestructuring5.ts] ////

//// [declarationEmitDestructuring5.ts]
function baz([, z, , ]) { }
function foo([, b, ]: [any, any]): void { }
function bar([z, , , ]) { }
function bar1([z, , , ] = [1, 3, 4, 6, 7]) { }
function bar2([,,z, , , ]) { }

//// [declarationEmitDestructuring5.js]
"use strict";
function baz([, z, ,]) { }
function foo([, b,]) { }
function bar([z, , ,]) { }
function bar1([z, , ,] = [1, 3, 4, 6, 7]) { }
function bar2([, , z, , ,]) { }


//// [declarationEmitDestructuring5.d.ts]
function baz([, z, ,]: [any, any, any?]): void;
function foo([, b,]: [any, any]): void;
function bar([z, , ,]: [any, any?, any?]): void;
function bar1([z, , ,]?: [number, number, number, number, number]): void;
function bar2([, , z, , ,]: [any, any, any, any?, any?]): void;


//// [DtsFileErrors]


declarationEmitDestructuring5.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitDestructuring5.d.ts (1 errors) ====
    function baz([, z, ,]: [any, any, any?]): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function foo([, b,]: [any, any]): void;
    function bar([z, , ,]: [any, any?, any?]): void;
    function bar1([z, , ,]?: [number, number, number, number, number]): void;
    function bar2([, , z, , ,]: [any, any, any, any?, any?]): void;
    