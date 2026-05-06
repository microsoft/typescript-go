//// [tests/cases/conformance/es6/destructuring/emptyAssignmentPatterns01_ES5.ts] ////

//// [emptyAssignmentPatterns01_ES5.ts]
var a: any;

({} = a);
([] = a);

var [,] = [1,2];

//// [emptyAssignmentPatterns01_ES5.js]
"use strict";
var a;
({} = a);
([] = a);
var [,] = [1, 2];


//// [emptyAssignmentPatterns01_ES5.d.ts]
var a: any;


//// [DtsFileErrors]


emptyAssignmentPatterns01_ES5.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== emptyAssignmentPatterns01_ES5.d.ts (1 errors) ====
    var a: any;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    