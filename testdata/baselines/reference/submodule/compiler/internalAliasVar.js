//// [tests/cases/compiler/internalAliasVar.ts] ////

//// [internalAliasVar.ts]
namespace a {
    export var x = 10;
}

namespace c {
    import b = a.x;
    export var bVal = b;
}


//// [internalAliasVar.js]
"use strict";
var a;
(function (a) {
    a.x = 10;
})(a || (a = {}));
var c;
(function (c) {
    var b = a.x;
    c.bVal = b;
})(c || (c = {}));


//// [internalAliasVar.d.ts]
namespace a {
    var x: number;
}
namespace c {
    var bVal: number;
}


//// [DtsFileErrors]


internalAliasVar.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== internalAliasVar.d.ts (1 errors) ====
    namespace a {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        var x: number;
    }
    namespace c {
        var bVal: number;
    }
    