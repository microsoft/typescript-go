//// [tests/cases/compiler/internalAliasClass.ts] ////

//// [internalAliasClass.ts]
namespace a {
    export class c {
    }
}

namespace c {
    import b = a.c;
    export var x: b = new b();
}

//// [internalAliasClass.js]
"use strict";
var a;
(function (a) {
    class c {
    }
    a.c = c;
})(a || (a = {}));
var c;
(function (c) {
    var b = a.c;
    c.x = new b();
})(c || (c = {}));


//// [internalAliasClass.d.ts]
namespace a {
    class c {
    }
}
namespace c {
    import b = a.c;
    var x: b;
}


//// [DtsFileErrors]


internalAliasClass.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== internalAliasClass.d.ts (1 errors) ====
    namespace a {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        class c {
        }
    }
    namespace c {
        import b = a.c;
        var x: b;
    }
    