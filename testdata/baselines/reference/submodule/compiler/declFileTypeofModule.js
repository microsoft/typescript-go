//// [tests/cases/compiler/declFileTypeofModule.ts] ////

//// [declFileTypeofModule.ts]
namespace m1 {
    export var c: string;
}
var m1_1 = m1;
var m1_2: typeof m1;

namespace m2 {
    export var d: typeof m2;
}

var m2_1 = m2;
var m2_2: typeof m2;

//// [declFileTypeofModule.js]
"use strict";
var m1;
(function (m1) {
})(m1 || (m1 = {}));
var m1_1 = m1;
var m1_2;
var m2;
(function (m2) {
})(m2 || (m2 = {}));
var m2_1 = m2;
var m2_2;


//// [declFileTypeofModule.d.ts]
namespace m1 {
    var c: string;
}
var m1_1: typeof m1;
var m1_2: typeof m1;
namespace m2 {
    var d: typeof m2;
}
var m2_1: typeof m2;
var m2_2: typeof m2;


//// [DtsFileErrors]


declFileTypeofModule.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileTypeofModule.d.ts (1 errors) ====
    namespace m1 {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        var c: string;
    }
    var m1_1: typeof m1;
    var m1_2: typeof m1;
    namespace m2 {
        var d: typeof m2;
    }
    var m2_1: typeof m2;
    var m2_2: typeof m2;
    