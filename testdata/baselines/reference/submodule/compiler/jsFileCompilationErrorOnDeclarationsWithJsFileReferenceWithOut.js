//// [tests/cases/compiler/jsFileCompilationErrorOnDeclarationsWithJsFileReferenceWithOut.ts] ////

//// [a.ts]
class c {
}

//// [b.ts]
/// <reference path="c.js"/>
function foo() {
}

//// [c.js]
function bar() {
}


//// [a.js]
class c {
}
//// [b.js]
function foo() {
}
//// [c.js]
function bar() {
}
