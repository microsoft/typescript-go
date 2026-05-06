//// [tests/cases/compiler/structuralTypeInDeclareFileForModule.ts] ////

//// [structuralTypeInDeclareFileForModule.ts]
namespace M { export var x; }
var m = M;

//// [structuralTypeInDeclareFileForModule.js]
"use strict";
var M;
(function (M) {
})(M || (M = {}));
var m = M;


//// [structuralTypeInDeclareFileForModule.d.ts]
namespace M {
    var x: any;
}
var m: typeof M;


//// [DtsFileErrors]


structuralTypeInDeclareFileForModule.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== structuralTypeInDeclareFileForModule.d.ts (1 errors) ====
    namespace M {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        var x: any;
    }
    var m: typeof M;
    