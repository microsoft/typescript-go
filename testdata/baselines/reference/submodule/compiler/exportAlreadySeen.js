//// [tests/cases/compiler/exportAlreadySeen.ts] ////

//// [exportAlreadySeen.ts]
module M {
    export export var x = 1;
    export export function f() { }

    export export module N {
        export export class C { }
        export export interface I { }
    }  
}

declare module A {
    export export var x;
    export export function f()

    export export module N {
        export export class C { }
        export export interface I { }
    }
}

//// [exportAlreadySeen.js]
var M;
(function (M) {
    M.x = 1;
    function f() { }
    M.f = f;
    let N;
    (function (N) {
        class C {
        }
        N.C = C;
    })(N = M.N || (M.N = {}));
})(M || (M = {}));
