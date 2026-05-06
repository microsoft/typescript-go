//// [tests/cases/compiler/genericClassesInModule.ts] ////

//// [genericClassesInModule.ts]
namespace Foo {

    export class B<T>{ }

    export class A { }
}

var a = new Foo.B<Foo.A>();

//// [genericClassesInModule.js]
"use strict";
var Foo;
(function (Foo) {
    class B {
    }
    Foo.B = B;
    class A {
    }
    Foo.A = A;
})(Foo || (Foo = {}));
var a = new Foo.B();


//// [genericClassesInModule.d.ts]
namespace Foo {
    class B<T> {
    }
    class A {
    }
}
var a: Foo.B<Foo.A>;


//// [DtsFileErrors]


genericClassesInModule.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== genericClassesInModule.d.ts (1 errors) ====
    namespace Foo {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        class B<T> {
        }
        class A {
        }
    }
    var a: Foo.B<Foo.A>;
    