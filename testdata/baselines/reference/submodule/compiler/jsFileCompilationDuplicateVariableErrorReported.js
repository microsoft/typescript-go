//// [tests/cases/compiler/jsFileCompilationDuplicateVariableErrorReported.ts] ////

//// [b.js]
var x = "hello";

//// [a.ts]
var x = 10; // Error reported


//// [a.js]
var x = 10;
//// [b.js]
var x = "hello";
