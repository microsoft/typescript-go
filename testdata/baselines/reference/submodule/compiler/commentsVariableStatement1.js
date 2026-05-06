//// [tests/cases/compiler/commentsVariableStatement1.ts] ////

//// [commentsVariableStatement1.ts]
/** Comment */
var v = 1;

//// [commentsVariableStatement1.js]
"use strict";
/** Comment */
var v = 1;


//// [commentsVariableStatement1.d.ts]
/** Comment */
var v: number;


//// [DtsFileErrors]


commentsVariableStatement1.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== commentsVariableStatement1.d.ts (1 errors) ====
    /** Comment */
    var v: number;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    