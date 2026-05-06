//// [tests/cases/conformance/es6/destructuring/emptyAssignmentPatterns03_ES6.ts] ////

//// [emptyAssignmentPatterns03_ES6.ts]
var a: any;

({} = {} = a);
([] = [] = a);

//// [emptyAssignmentPatterns03_ES6.js]
"use strict";
var a;
({} = {} = a);
([] = [] = a);


//// [emptyAssignmentPatterns03_ES6.d.ts]
var a: any;


//// [DtsFileErrors]


emptyAssignmentPatterns03_ES6.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== emptyAssignmentPatterns03_ES6.d.ts (1 errors) ====
    var a: any;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    