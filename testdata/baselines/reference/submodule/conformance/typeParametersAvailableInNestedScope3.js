//// [tests/cases/conformance/types/typeParameters/typeParameterLists/typeParametersAvailableInNestedScope3.ts] ////

//// [typeParametersAvailableInNestedScope3.ts]
function foo<T>(v: T) {
    function a<T>(a: T) { return a; }
    function b(): T { return v; }

    function c<T>(v: T) {
        function a<T>(a: T) { return a; }
        function b(): T { return v; }
        return { a, b };
    }

    return { a, b, c };
}


//// [typeParametersAvailableInNestedScope3.js]
"use strict";
function foo(v) {
    function a(a) { return a; }
    function b() { return v; }
    function c(v) {
        function a(a) { return a; }
        function b() { return v; }
        return { a, b };
    }
    return { a, b, c };
}


//// [typeParametersAvailableInNestedScope3.d.ts]
function foo<T>(v: T): {
    a: <T_1>(a: T_1) => T_1;
    b: () => T;
    c: <T_1>(v: T_1) => {
        a: <T_2>(a: T_2) => T_2;
        b: () => T_1;
    };
};


//// [DtsFileErrors]


typeParametersAvailableInNestedScope3.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== typeParametersAvailableInNestedScope3.d.ts (1 errors) ====
    function foo<T>(v: T): {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        a: <T_1>(a: T_1) => T_1;
        b: () => T;
        c: <T_1>(v: T_1) => {
            a: <T_2>(a: T_2) => T_2;
            b: () => T_1;
        };
    };
    