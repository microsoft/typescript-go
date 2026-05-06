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

//// [a.js]
"use strict";
class c {
}
//// [b.js]
"use strict";
/// <reference path="c.js"/>
// b.d.ts should have c.d.ts as the reference path
function foo() {
}


//// [a.d.ts]
class c {
}
//// [c.d.ts]
function bar(): void;
//// [b.d.ts]
function foo(): void;
