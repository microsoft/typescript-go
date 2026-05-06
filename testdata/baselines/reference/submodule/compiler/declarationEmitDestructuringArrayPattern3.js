//// [tests/cases/compiler/declarationEmitDestructuringArrayPattern3.ts] ////

//// [declarationEmitDestructuringArrayPattern3.ts]
namespace M {
    export var [a, b] = [1, 2];
}

//// [declarationEmitDestructuringArrayPattern3.js]
"use strict";
var M;
(function (M) {
    var _a;
    _a = [1, 2], M.a = _a[0], M.b = _a[1];
})(M || (M = {}));


//// [declarationEmitDestructuringArrayPattern3.d.ts]
namespace M {
    var a: number, b: number;
}


//// [DtsFileErrors]


declarationEmitDestructuringArrayPattern3.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitDestructuringArrayPattern3.d.ts (1 errors) ====
    namespace M {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        var a: number, b: number;
    }
    