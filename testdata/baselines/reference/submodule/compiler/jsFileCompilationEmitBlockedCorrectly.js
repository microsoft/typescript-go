//// [tests/cases/compiler/jsFileCompilationEmitBlockedCorrectly.ts] ////

//// [a.ts]
class c {
}

//// [b.ts]
// this should be emitted
class d {
}

//// [a.js]
function foo() {
}


//// [a.js]
class c {
}
//// [b.js]
class d {
}
