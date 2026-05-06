//// [tests/cases/compiler/declFileObjectLiteralWithOnlySetter.ts] ////

//// [declFileObjectLiteralWithOnlySetter.ts]
function /*1*/makePoint(x: number) { 
    return {
        b: 10,
        set x(a: number) { this.b = a; }
    };
};
var /*3*/point = makePoint(2);
point./*2*/x = 30;

//// [declFileObjectLiteralWithOnlySetter.js]
"use strict";
function makePoint(x) {
    return {
        b: 10,
        set x(a) { this.b = a; }
    };
}
;
var /*3*/ point = makePoint(2);
point. /*2*/x = 30;


//// [declFileObjectLiteralWithOnlySetter.d.ts]
function makePoint(x: number): {
    b: number;
    x: number;
};
var /*3*/ point: {
    b: number;
    x: number;
};


//// [DtsFileErrors]


declFileObjectLiteralWithOnlySetter.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileObjectLiteralWithOnlySetter.d.ts (1 errors) ====
    function makePoint(x: number): {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        b: number;
        x: number;
    };
    var /*3*/ point: {
        b: number;
        x: number;
    };
    