//// [tests/cases/compiler/declFileObjectLiteralWithAccessors.ts] ////

//// [declFileObjectLiteralWithAccessors.ts]
function /*1*/makePoint(x: number) { 
    return {
        b: 10,
        get x() { return x; },
        set x(a: number) { this.b = a; }
    };
};
var /*4*/point = makePoint(2);
var /*2*/x = point.x;
point./*3*/x = 30;

//// [declFileObjectLiteralWithAccessors.js]
"use strict";
function makePoint(x) {
    return {
        b: 10,
        get x() { return x; },
        set x(a) { this.b = a; }
    };
}
;
var /*4*/ point = makePoint(2);
var /*2*/ x = point.x;
point. /*3*/x = 30;


//// [declFileObjectLiteralWithAccessors.d.ts]
function makePoint(x: number): {
    b: number;
    x: number;
};
var /*4*/ point: {
    b: number;
    x: number;
};
var /*2*/ x: number;


//// [DtsFileErrors]


declFileObjectLiteralWithAccessors.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileObjectLiteralWithAccessors.d.ts (1 errors) ====
    function makePoint(x: number): {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        b: number;
        x: number;
    };
    var /*4*/ point: {
        b: number;
        x: number;
    };
    var /*2*/ x: number;
    