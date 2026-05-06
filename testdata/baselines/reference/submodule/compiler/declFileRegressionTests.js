//// [tests/cases/compiler/declFileRegressionTests.ts] ////

//// [declFileRegressionTests.ts]
// 'null' not converted to 'any' in d.ts
// function types not piped through correctly
var n = { w: null, x: '', y: () => { }, z: 32 };



//// [declFileRegressionTests.js]
"use strict";
// 'null' not converted to 'any' in d.ts
// function types not piped through correctly
var n = { w: null, x: '', y: () => { }, z: 32 };


//// [declFileRegressionTests.d.ts]
var n: {
    w: any;
    x: string;
    y: () => void;
    z: number;
};


//// [DtsFileErrors]


declFileRegressionTests.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileRegressionTests.d.ts (1 errors) ====
    var n: {
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        w: any;
        x: string;
        y: () => void;
        z: number;
    };
    