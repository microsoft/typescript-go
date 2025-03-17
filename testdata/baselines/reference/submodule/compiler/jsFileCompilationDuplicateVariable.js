//// [tests/cases/compiler/jsFileCompilationDuplicateVariable.ts] ////

//// [a.ts]
var x = 10;

//// [b.js]
var x = "hello"; // Error is recorded here, but suppressed because the js file isn't checked


//// [b.js]
var x = "hello";
//// [a.js]
var x = 10;
