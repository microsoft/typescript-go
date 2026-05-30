//// [tests/cases/compiler/loneSurrogateStringLiteralEmit.ts] ////

//// [loneSurrogateStringLiteralEmit.ts]
enum E {
  "\uD800lone" = 1,
}

const enum S {
  C = "\uD800",
}

const c = S.C;


//// [loneSurrogateStringLiteralEmit.js]
"use strict";
var E;
(function (E) {
    E[E["\uD800lone"] = 1] = "\uD800lone";
})(E || (E = {}));
const c = "\uD800" /* S.C */;
