//// [tests/cases/compiler/ambientConstLiterals.ts] ////

//// [ambientConstLiterals.ts]
function f<T>(x: T): T {
    return x;
}

enum E { A, B, C, "non identifier" }

const c1 = "abc";
const c2 = 123;
const c3 = c1;
const c4 = c2;
const c5 = f(123);
const c6 = f(-123);
const c7 = true;
const c8 = E.A;
const c8b = E["non identifier"];
const c9 = { x: "abc" };
const c10 = [123];
const c11 = "abc" + "def";
const c12 = 123 + 456;
const c13 = Math.random() > 0.5 ? "abc" : "def";
const c14 = Math.random() > 0.5 ? 123 : 456;

//// [ambientConstLiterals.js]
"use strict";
function f(x) {
    return x;
}
var E;
(function (E) {
    E[E["A"] = 0] = "A";
    E[E["B"] = 1] = "B";
    E[E["C"] = 2] = "C";
    E[E["non identifier"] = 3] = "non identifier";
})(E || (E = {}));
const c1 = "abc";
const c2 = 123;
const c3 = c1;
const c4 = c2;
const c5 = f(123);
const c6 = f(-123);
const c7 = true;
const c8 = E.A;
const c8b = E["non identifier"];
const c9 = { x: "abc" };
const c10 = [123];
const c11 = "abc" + "def";
const c12 = 123 + 456;
const c13 = Math.random() > 0.5 ? "abc" : "def";
const c14 = Math.random() > 0.5 ? 123 : 456;


//// [ambientConstLiterals.d.ts]
function f<T>(x: T): T;
enum E {
    A = 0,
    B = 1,
    C = 2,
    "non identifier" = 3
}
const c1 = "abc";
const c2 = 123;
const c3 = "abc";
const c4 = 123;
const c5 = 123;
const c6 = -123;
const c7 = true;
const c8 = E.A;
const c8b = E["non identifier"];
const c9: {
    x: string;
};
const c10: number[];
const c11: string;
const c12: number;
const c13: string;
const c14: number;


//// [DtsFileErrors]


ambientConstLiterals.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== ambientConstLiterals.d.ts (1 errors) ====
    function f<T>(x: T): T;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    enum E {
        A = 0,
        B = 1,
        C = 2,
        "non identifier" = 3
    }
    const c1 = "abc";
    const c2 = 123;
    const c3 = "abc";
    const c4 = 123;
    const c5 = 123;
    const c6 = -123;
    const c7 = true;
    const c8 = E.A;
    const c8b = E["non identifier"];
    const c9: {
        x: string;
    };
    const c10: number[];
    const c11: string;
    const c12: number;
    const c13: string;
    const c14: number;
    