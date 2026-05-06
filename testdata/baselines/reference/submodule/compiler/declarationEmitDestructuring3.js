//// [tests/cases/compiler/declarationEmitDestructuring3.ts] ////

//// [declarationEmitDestructuring3.ts]
function bar([x, z, ...w]) { }
function foo([x, ...y] = [1, "string", true]) { }



//// [declarationEmitDestructuring3.js]
"use strict";
function bar([x, z, ...w]) { }
function foo([x, ...y] = [1, "string", true]) { }


//// [declarationEmitDestructuring3.d.ts]
function bar([x, z, ...w]: [any, any, ...any[]]): void;
function foo([x, ...y]?: [number, string, boolean]): void;


//// [DtsFileErrors]


declarationEmitDestructuring3.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitDestructuring3.d.ts (1 errors) ====
    function bar([x, z, ...w]: [any, any, ...any[]]): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function foo([x, ...y]?: [number, string, boolean]): void;
    