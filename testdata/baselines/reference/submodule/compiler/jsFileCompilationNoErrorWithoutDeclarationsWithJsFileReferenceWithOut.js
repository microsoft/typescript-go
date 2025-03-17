//// [tests/cases/compiler/jsFileCompilationNoErrorWithoutDeclarationsWithJsFileReferenceWithOut.ts] ////

//// [a.ts]
class c {
}

//// [b.ts]
/// <reference path="c.js"/>
//no  error on above reference since not emitting declarations
function foo() {
}

//// [c.js]
function bar() {
}


//// [b.js]
function foo() {
}
//// [c.js]
function bar() {
}
//// [a.js]
class c {
}
