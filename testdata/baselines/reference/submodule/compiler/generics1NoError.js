//// [tests/cases/compiler/generics1NoError.ts] ////

//// [generics1NoError.ts]
interface A { a: string; }
interface B extends A { b: string; }
interface C extends B { c: string; }
interface G<T, U extends B> {
    x: T;
    y: U;
}
var v1: G<A, C>;               // Ok
var v2: G<{ a: string }, C>;   // Ok, equivalent to G<A, C>
var v4: G<G<A, B>, C>;         // Ok

//// [generics1NoError.js]
"use strict";
var v1; // Ok
var v2; // Ok, equivalent to G<A, C>
var v4; // Ok


//// [generics1NoError.d.ts]
interface A {
    a: string;
}
interface B extends A {
    b: string;
}
interface C extends B {
    c: string;
}
interface G<T, U extends B> {
    x: T;
    y: U;
}
var v1: G<A, C>;
var v2: G<{
    a: string;
}, C>;
var v4: G<G<A, B>, C>;


//// [DtsFileErrors]


generics1NoError.d.ts(14,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== generics1NoError.d.ts (1 errors) ====
    interface A {
        a: string;
    }
    interface B extends A {
        b: string;
    }
    interface C extends B {
        c: string;
    }
    interface G<T, U extends B> {
        x: T;
        y: U;
    }
    var v1: G<A, C>;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    var v2: G<{
        a: string;
    }, C>;
    var v4: G<G<A, B>, C>;
    