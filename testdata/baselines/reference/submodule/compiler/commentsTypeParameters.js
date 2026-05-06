//// [tests/cases/compiler/commentsTypeParameters.ts] ////

//// [commentsTypeParameters.ts]
class C</**docComment for type parameter*/ T> {
    method</**docComment of method type parameter */ U extends T>(a: U) {
    }
    static staticmethod</**docComment of method type parameter */ U>(a: U) {
    }

    private privatemethod</**docComment of method type parameter */ U extends T>(a: U) {
    }
    private static privatestaticmethod</**docComment of method type parameter */ U>(a: U) {
    }
}

function compare</**type*/T>(a: T, b: T) {
    return a === b;
}

//// [commentsTypeParameters.js]
"use strict";
class C {
    method(a) {
    }
    static staticmethod(a) {
    }
    privatemethod(a) {
    }
    static privatestaticmethod(a) {
    }
}
function compare(a, b) {
    return a === b;
}


//// [commentsTypeParameters.d.ts]
class C</**docComment for type parameter*/ T> {
    method</**docComment of method type parameter */ U extends T>(a: U): void;
    static staticmethod</**docComment of method type parameter */ U>(a: U): void;
    private privatemethod;
    private static privatestaticmethod;
}
function compare</**type*/ T>(a: T, b: T): boolean;


//// [DtsFileErrors]


commentsTypeParameters.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== commentsTypeParameters.d.ts (1 errors) ====
    class C</**docComment for type parameter*/ T> {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        method</**docComment of method type parameter */ U extends T>(a: U): void;
        static staticmethod</**docComment of method type parameter */ U>(a: U): void;
        private privatemethod;
        private static privatestaticmethod;
    }
    function compare</**type*/ T>(a: T, b: T): boolean;
    