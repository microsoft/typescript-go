//// [tests/cases/compiler/declFileObjectLiteralWithOnlyGetter.ts] ////

//// [declFileObjectLiteralWithOnlyGetter.ts]
function /*1*/makePoint(x: number) { 
    return {
        get x() { return x; },
    };
};
var /*4*/point = makePoint(2);
var /*2*/x = point./*3*/x;


//// [declFileObjectLiteralWithOnlyGetter.js]
"use strict";
function makePoint(x) {
    return {
        get x() { return x; },
    };
}
;
var /*4*/ point = makePoint(2);
var /*2*/ x = point. /*3*/x;


//// [declFileObjectLiteralWithOnlyGetter.d.ts]
function makePoint(x: number): {
    readonly x: number;
};
var /*4*/ point: {
    readonly x: number;
};
var /*2*/ x: number;


//// [DtsFileErrors]


declFileObjectLiteralWithOnlyGetter.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileObjectLiteralWithOnlyGetter.d.ts (1 errors) ====
    function makePoint(x: number): {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        readonly x: number;
    };
    var /*4*/ point: {
        readonly x: number;
    };
    var /*2*/ x: number;
    