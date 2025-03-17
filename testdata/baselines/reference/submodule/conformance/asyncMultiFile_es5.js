//// [tests/cases/conformance/async/es5/asyncMultiFile_es5.ts] ////

//// [a.ts]
async function f() {}
//// [b.ts]
function g() { }

//// [b.js]
function g() { }
//// [a.js]
async function f() { }
