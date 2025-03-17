//// [tests/cases/conformance/async/es6/asyncMultiFile_es6.ts] ////

//// [a.ts]
async function f() {}
//// [b.ts]
function g() { }

//// [b.js]
function g() { }
//// [a.js]
async function f() { }
