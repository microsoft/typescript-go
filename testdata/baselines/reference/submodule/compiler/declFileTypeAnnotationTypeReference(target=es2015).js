//// [tests/cases/compiler/declFileTypeAnnotationTypeReference.ts] ////

//// [declFileTypeAnnotationTypeReference.ts]
class c {
}
namespace m {
    export class c {
    }
    export class g<T> {
    }
}
class g<T> {
}

// Just the name
function foo(): c {
    return new c();
}
function foo2() {
    return new c();
}

// Qualified name
function foo3(): m.c {
    return new m.c();
}
function foo4() {
    return new m.c();
}

// Just the name with type arguments
function foo5(): g<string> {
    return new g<string>();
}
function foo6() {
    return new g<string>();
}

// Qualified name with type arguments
function foo7(): m.g<number> {
    return new m.g<number>();
}
function foo8() {
    return new m.g<number>();
}

//// [declFileTypeAnnotationTypeReference.js]
"use strict";
class c {
}
var m;
(function (m) {
    class c {
    }
    m.c = c;
    class g {
    }
    m.g = g;
})(m || (m = {}));
class g {
}
// Just the name
function foo() {
    return new c();
}
function foo2() {
    return new c();
}
// Qualified name
function foo3() {
    return new m.c();
}
function foo4() {
    return new m.c();
}
// Just the name with type arguments
function foo5() {
    return new g();
}
function foo6() {
    return new g();
}
// Qualified name with type arguments
function foo7() {
    return new m.g();
}
function foo8() {
    return new m.g();
}


//// [declFileTypeAnnotationTypeReference.d.ts]
class c {
}
namespace m {
    class c {
    }
    class g<T> {
    }
}
class g<T> {
}
function foo(): c;
function foo2(): c;
function foo3(): m.c;
function foo4(): m.c;
function foo5(): g<string>;
function foo6(): g<string>;
function foo7(): m.g<number>;
function foo8(): m.g<number>;


//// [DtsFileErrors]


declFileTypeAnnotationTypeReference.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileTypeAnnotationTypeReference.d.ts (1 errors) ====
    class c {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    namespace m {
        class c {
        }
        class g<T> {
        }
    }
    class g<T> {
    }
    function foo(): c;
    function foo2(): c;
    function foo3(): m.c;
    function foo4(): m.c;
    function foo5(): g<string>;
    function foo6(): g<string>;
    function foo7(): m.g<number>;
    function foo8(): m.g<number>;
    