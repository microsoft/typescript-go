//// [tests/cases/compiler/moduleSymbolMerging.ts] ////

//// [A.ts]
module A { export interface I {} }

//// [B.ts]
///<reference path="A.ts" preserve="true" />
module A { ; }
module B {
	export function f(): A.I { return null; }
}



//// [B.js]
var A;
(function (A) {
    ;
})(A || (A = {}));
var B;
(function (B) {
    function f() { return null; }
    B.f = f;
})(B || (B = {}));
//// [A.js]
