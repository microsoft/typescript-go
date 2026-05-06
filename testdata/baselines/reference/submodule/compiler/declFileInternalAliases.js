//// [tests/cases/compiler/declFileInternalAliases.ts] ////

//// [declFileInternalAliases.ts]
namespace m {
    export class c {
    }
}
namespace m1 {
    import x = m.c;
    export var d = new x(); // emit the type as m.c
}
namespace m2 {
    export import x = m.c;
    export var d = new x(); // emit the type as x
}

//// [declFileInternalAliases.js]
"use strict";
var m;
(function (m) {
    class c {
    }
    m.c = c;
})(m || (m = {}));
var m1;
(function (m1) {
    var x = m.c;
    m1.d = new x(); // emit the type as m.c
})(m1 || (m1 = {}));
var m2;
(function (m2) {
    m2.x = m.c;
    m2.d = new m2.x(); // emit the type as x
})(m2 || (m2 = {}));


//// [declFileInternalAliases.d.ts]
namespace m {
    class c {
    }
}
namespace m1 {
    import x = m.c;
    var d: x;
}
namespace m2 {
    export import x = m.c;
    var d: x;
}


//// [DtsFileErrors]


declFileInternalAliases.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileInternalAliases.d.ts (1 errors) ====
    namespace m {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        class c {
        }
    }
    namespace m1 {
        import x = m.c;
        var d: x;
    }
    namespace m2 {
        export import x = m.c;
        var d: x;
    }
    