//// [tests/cases/compiler/jsFileCompilationLetDeclarationOrder.ts] ////

//// [b.js]
let a = 10;
b = 30;

//// [a.ts]
let b = 30;
a = 10;


//// [a.js]
let b = 30;
a = 10;
//// [b.js]
let a = 10;
b = 30;
