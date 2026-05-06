//// [tests/cases/compiler/genericClassImplementingGenericInterfaceFromAnotherModule.ts] ////

//// [genericClassImplementingGenericInterfaceFromAnotherModule.ts]
namespace foo {
    export interface IFoo<T> { }
}
namespace bar {
    export class Foo<T> implements foo.IFoo<T> { }
}


//// [genericClassImplementingGenericInterfaceFromAnotherModule.js]
"use strict";
var bar;
(function (bar) {
    class Foo {
    }
    bar.Foo = Foo;
})(bar || (bar = {}));


//// [genericClassImplementingGenericInterfaceFromAnotherModule.d.ts]
namespace foo {
    interface IFoo<T> {
    }
}
namespace bar {
    class Foo<T> implements foo.IFoo<T> {
    }
}


//// [DtsFileErrors]


genericClassImplementingGenericInterfaceFromAnotherModule.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== genericClassImplementingGenericInterfaceFromAnotherModule.d.ts (1 errors) ====
    namespace foo {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        interface IFoo<T> {
        }
    }
    namespace bar {
        class Foo<T> implements foo.IFoo<T> {
        }
    }
    