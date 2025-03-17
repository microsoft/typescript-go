//// [tests/cases/compiler/alwaysStrictModule2.ts] ////

//// [a.ts]
module M {
    export function f() {
        var arguments = [];
    }
}

//// [b.ts]
module M {
    export function f2() {
        var arguments = [];
    }
}

//// [b.js]
var M;
(function (M) {
    function f2() {
        var arguments = [];
    }
    M.f2 = f2;
})(M || (M = {}));
//// [a.js]
var M;
(function (M) {
    function f() {
        var arguments = [];
    }
    M.f = f;
})(M || (M = {}));
