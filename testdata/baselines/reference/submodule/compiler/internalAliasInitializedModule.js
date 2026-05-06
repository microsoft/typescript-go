//// [tests/cases/compiler/internalAliasInitializedModule.ts] ////

//// [internalAliasInitializedModule.ts]
namespace a {
    export namespace b {
        export class c {
        }
    }
}

namespace c {
    import b = a.b;
    export var x: b.c = new b.c();
}

//// [internalAliasInitializedModule.js]
"use strict";
var a;
(function (a) {
    let b;
    (function (b) {
        class c {
        }
        b.c = c;
    })(b = a.b || (a.b = {}));
})(a || (a = {}));
var c;
(function (c) {
    var b = a.b;
    c.x = new b.c();
})(c || (c = {}));


//// [internalAliasInitializedModule.d.ts]
namespace a {
    namespace b {
        class c {
        }
    }
}
namespace c {
    import b = a.b;
    var x: b.c;
}


//// [DtsFileErrors]


internalAliasInitializedModule.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== internalAliasInitializedModule.d.ts (1 errors) ====
    namespace a {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        namespace b {
            class c {
            }
        }
    }
    namespace c {
        import b = a.b;
        var x: b.c;
    }
    