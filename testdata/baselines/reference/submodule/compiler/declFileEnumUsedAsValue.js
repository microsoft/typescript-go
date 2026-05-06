//// [tests/cases/compiler/declFileEnumUsedAsValue.ts] ////

//// [declFileEnumUsedAsValue.ts]
enum e {
    a,
    b,
    c
}
var x = e;

//// [declFileEnumUsedAsValue.js]
"use strict";
var e;
(function (e) {
    e[e["a"] = 0] = "a";
    e[e["b"] = 1] = "b";
    e[e["c"] = 2] = "c";
})(e || (e = {}));
var x = e;


//// [declFileEnumUsedAsValue.d.ts]
enum e {
    a = 0,
    b = 1,
    c = 2
}
var x: typeof e;


//// [DtsFileErrors]


declFileEnumUsedAsValue.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileEnumUsedAsValue.d.ts (1 errors) ====
    enum e {
    ~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        a = 0,
        b = 1,
        c = 2
    }
    var x: typeof e;
    