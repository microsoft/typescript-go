//// [tests/cases/conformance/es6/destructuring/emptyAssignmentPatterns04_ES5.ts] ////

//// [emptyAssignmentPatterns04_ES5.ts]
var a: any;
let x, y, z, a1, a2, a3;

({ x, y, z } = {} = a);
([ a1, a2, a3] = [] = a);

//// [emptyAssignmentPatterns04_ES5.js]
"use strict";
var a;
let x, y, z, a1, a2, a3;
({ x, y, z } = {} = a);
([a1, a2, a3] = [] = a);


//// [emptyAssignmentPatterns04_ES5.d.ts]
var a: any;
let x: any, y: any, z: any, a1: any, a2: any, a3: any;


//// [DtsFileErrors]


emptyAssignmentPatterns04_ES5.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== emptyAssignmentPatterns04_ES5.d.ts (1 errors) ====
    var a: any;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    let x: any, y: any, z: any, a1: any, a2: any, a3: any;
    