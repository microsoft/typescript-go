//// [tests/cases/compiler/internalAliasInterface.ts] ////

//// [internalAliasInterface.ts]
namespace a {
    export interface I {
    }
}

namespace c {
    import b = a.I;
    export var x: b;
}


//// [internalAliasInterface.js]
"use strict";
var c;
(function (c) {
})(c || (c = {}));


//// [internalAliasInterface.d.ts]
namespace a {
    interface I {
    }
}
namespace c {
    import b = a.I;
    var x: b;
}


//// [DtsFileErrors]


internalAliasInterface.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== internalAliasInterface.d.ts (1 errors) ====
    namespace a {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        interface I {
        }
    }
    namespace c {
        import b = a.I;
        var x: b;
    }
    