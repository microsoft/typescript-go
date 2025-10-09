//// [tests/cases/compiler/alwaysStrictModule.ts] ////

//// [alwaysStrictModule.ts]
namespace M {
    export function f() {
        var arguments = [];
    }
}

//// [alwaysStrictModule.js]
var M;
(function (M) {
    function f() {
        var arguments = [];
    }
    M.f = f;
})(M || (M = {}));
