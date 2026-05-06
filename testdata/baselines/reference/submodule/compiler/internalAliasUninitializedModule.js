//// [tests/cases/compiler/internalAliasUninitializedModule.ts] ////

//// [internalAliasUninitializedModule.ts]
namespace a {
    export namespace b {
        export interface I {
            foo();
        }
    }
}

namespace c {
    import b = a.b;
    export var x: b.I;
    x.foo();
}

//// [internalAliasUninitializedModule.js]
"use strict";
var c;
(function (c) {
    c.x.foo();
})(c || (c = {}));


//// [internalAliasUninitializedModule.d.ts]
namespace a {
    namespace b {
        interface I {
            foo(): any;
        }
    }
}
namespace c {
    import b = a.b;
    var x: b.I;
}


//// [DtsFileErrors]


internalAliasUninitializedModule.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== internalAliasUninitializedModule.d.ts (1 errors) ====
    namespace a {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        namespace b {
            interface I {
                foo(): any;
            }
        }
    }
    namespace c {
        import b = a.b;
        var x: b.I;
    }
    