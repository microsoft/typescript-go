//// [tests/cases/compiler/declarationEmitInvalidReference.ts] ////

//// [declarationEmitInvalidReference.ts]
/// <reference path="invalid.ts" />
var x = 0;

//// [declarationEmitInvalidReference.js]
"use strict";
/// <reference path="invalid.ts" />
var x = 0;


//// [declarationEmitInvalidReference.d.ts]
var x: number;


//// [DtsFileErrors]


declarationEmitInvalidReference.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitInvalidReference.d.ts (1 errors) ====
    var x: number;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    