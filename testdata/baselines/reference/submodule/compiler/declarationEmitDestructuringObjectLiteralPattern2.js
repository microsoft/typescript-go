//// [tests/cases/compiler/declarationEmitDestructuringObjectLiteralPattern2.ts] ////

//// [declarationEmitDestructuringObjectLiteralPattern2.ts]
var { a: x11, b: { a: y11, b: { a: z11 }}} = { a: 1, b: { a: "hello", b: { a: true } } };

function f15() {
    var a4 = "hello";
    var b4 = 1;
    var c4 = true;
    return { a4, b4, c4 };
}
var { a4, b4, c4 } = f15();

namespace m {
    export var { a4, b4, c4 } = f15();
}

//// [declarationEmitDestructuringObjectLiteralPattern2.js]
"use strict";
var { a: x11, b: { a: y11, b: { a: z11 } } } = { a: 1, b: { a: "hello", b: { a: true } } };
function f15() {
    var a4 = "hello";
    var b4 = 1;
    var c4 = true;
    return { a4, b4, c4 };
}
var { a4, b4, c4 } = f15();
var m;
(function (m) {
    var _a;
    _a = f15(), m.a4 = _a.a4, m.b4 = _a.b4, m.c4 = _a.c4;
})(m || (m = {}));


//// [declarationEmitDestructuringObjectLiteralPattern2.d.ts]
var x11: number, y11: string, z11: boolean;
function f15(): {
    a4: string;
    b4: number;
    c4: boolean;
};
var a4: string, b4: number, c4: boolean;
namespace m {
    var a4: string, b4: number, c4: boolean;
}


//// [DtsFileErrors]


declarationEmitDestructuringObjectLiteralPattern2.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitDestructuringObjectLiteralPattern2.d.ts (1 errors) ====
    var x11: number, y11: string, z11: boolean;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function f15(): {
        a4: string;
        b4: number;
        c4: boolean;
    };
    var a4: string, b4: number, c4: boolean;
    namespace m {
        var a4: string, b4: number, c4: boolean;
    }
    