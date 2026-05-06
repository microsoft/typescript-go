//// [tests/cases/compiler/internalAliasWithDottedNameEmit.ts] ////

//// [internalAliasWithDottedNameEmit.ts]
namespace a.b.c {
      export var d;
}
namespace a.e.f {
      import g = b.c;
}


//// [internalAliasWithDottedNameEmit.js]
"use strict";
var a;
(function (a) {
    var b;
    (function (b) {
        var c;
        (function (c) {
        })(c = b.c || (b.c = {}));
    })(b = a.b || (a.b = {}));
})(a || (a = {}));


//// [internalAliasWithDottedNameEmit.d.ts]
namespace a.b.c {
    var d: any;
}
namespace a.e.f {
}


//// [DtsFileErrors]


internalAliasWithDottedNameEmit.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== internalAliasWithDottedNameEmit.d.ts (1 errors) ====
    namespace a.b.c {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        var d: any;
    }
    namespace a.e.f {
    }
    