//// [tests/cases/compiler/jsFileCompilationErrorOnDeclarationsWithJsFileReferenceWithNoOut.ts] ////

//// [a.ts]
class c {
}

//// [b.ts]
/// <reference path="c.js"/>
// b.d.ts should have c.d.ts as the reference path
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
