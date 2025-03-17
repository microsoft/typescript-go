//// [tests/cases/compiler/jsFileCompilationDuplicateFunctionImplementationFileOrderReversed.ts] ////

//// [a.ts]
function foo() {
    return 30;
}

//// [b.js]
function foo() {
    return 10;
}



//// [b.js]
function foo() {
    return 10;
}
//// [a.js]
function foo() {
    return 30;
}
