//// [tests/cases/compiler/jsFileCompilationEmitTrippleSlashReference.ts] ////

//// [a.ts]
class c {
}

//// [b.js]
/// <reference path="c.js"/>
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
