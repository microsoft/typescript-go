//// [tests/cases/conformance/es6/destructuring/emptyAssignmentPatterns01_ES5iterable.ts] ////

//// [emptyAssignmentPatterns01_ES5iterable.ts]
var a: any;

({} = a);
([] = a);

//// [emptyAssignmentPatterns01_ES5iterable.js]
"use strict";
var a;
({} = a);
([] = a);


//// [emptyAssignmentPatterns01_ES5iterable.d.ts]
var a: any;
