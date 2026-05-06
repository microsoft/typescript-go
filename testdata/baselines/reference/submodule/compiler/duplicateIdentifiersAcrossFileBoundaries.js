//// [tests/cases/compiler/duplicateIdentifiersAcrossFileBoundaries.ts] ////

//// [file1.ts]
interface I { }
class C1 { }
class C2 { }
function f() { }
var v = 3;

class Foo {
    static x: number;
}

namespace N {
    export namespace F {
        var t;
    }
}

//// [file2.ts]
class I { } // error -- cannot merge interface with non-ambient class
interface C1 { } // error -- cannot merge interface with non-ambient class
function C2() { } // error -- cannot merge function with non-ambient class
class f { } // error -- cannot merge function with non-ambient class
var v = 3;

namespace Foo {
    export var x: number; // error for redeclaring var in a different parent
}

declare namespace N {
    export function F(); // no error because function is ambient
}


//// [file1.js]
"use strict";
class C1 {
}
class C2 {
}
function f() { }
var v = 3;
class Foo {
}
var N;
(function (N) {
    let F;
    (function (F) {
        var t;
    })(F = N.F || (N.F = {}));
})(N || (N = {}));
//// [file2.js]
"use strict";
class I {
} // error -- cannot merge interface with non-ambient class
function C2() { } // error -- cannot merge function with non-ambient class
class f {
} // error -- cannot merge function with non-ambient class
var v = 3;
var Foo;
(function (Foo) {
})(Foo || (Foo = {}));


//// [file1.d.ts]
interface I {
}
class C1 {
}
class C2 {
}
function f(): void;
var v: number;
class Foo {
    static x: number;
}
namespace N {
    namespace F {
    }
}
//// [file2.d.ts]
class I {
}
interface C1 {
}
function C2(): void;
class f {
}
var v: number;
namespace Foo {
    var x: number;
}
namespace N {
    function F(): any;
}
