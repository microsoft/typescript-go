//// [tests/cases/compiler/jsFileCompilationLetBeingRenamed.ts] ////

//// [a.js]
function foo(a) {
    for (let a = 0; a < 10; a++) {
        // do something
    }
}


//// [a.js]
function foo(a) {
    for (let a = 0; a < 10; a++) {
    }
}
