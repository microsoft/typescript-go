//// [tests/cases/compiler/instantiatedTypeAliasDisplay.ts] ////

//// [instantiatedTypeAliasDisplay.ts]
// Repros from #12066

interface X<A> {
    a: A;
}
interface Y<B> {
    b: B;
}
type Z<A, B> = X<A> | Y<B>;

declare function f1<A>(): Z<A, number>;
declare function f2<A, B, C, D, E>(a: A, b: B, c: C, d: D): Z<A, string[]>;

const x1 = f1<string>();  // Z<string, number>
const x2 = f2({}, {}, {}, {});  // Z<{}, string[]>

//// [instantiatedTypeAliasDisplay.js]
"use strict";
// Repros from #12066
const x1 = f1(); // Z<string, number>
const x2 = f2({}, {}, {}, {}); // Z<{}, string[]>


//// [instantiatedTypeAliasDisplay.d.ts]
interface X<A> {
    a: A;
}
interface Y<B> {
    b: B;
}
type Z<A, B> = X<A> | Y<B>;
function f1<A>(): Z<A, number>;
function f2<A, B, C, D, E>(a: A, b: B, c: C, d: D): Z<A, string[]>;
const x1: Z<string, number>;
const x2: Z<{}, string[]>;


//// [DtsFileErrors]


instantiatedTypeAliasDisplay.d.ts(8,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== instantiatedTypeAliasDisplay.d.ts (1 errors) ====
    interface X<A> {
        a: A;
    }
    interface Y<B> {
        b: B;
    }
    type Z<A, B> = X<A> | Y<B>;
    function f1<A>(): Z<A, number>;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function f2<A, B, C, D, E>(a: A, b: B, c: C, d: D): Z<A, string[]>;
    const x1: Z<string, number>;
    const x2: Z<{}, string[]>;
    