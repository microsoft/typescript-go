//// [tests/cases/compiler/jsFileCompilationDuplicateFunctionImplementation.ts] ////

//// [b.js]
function foo() {
    return 10;
}
//// [a.ts]
function foo() {
    return 30;
}



//// [a.js]
function foo() {
    return 30;
}
//// [b.js]
function foo() {
    return 10;
}
