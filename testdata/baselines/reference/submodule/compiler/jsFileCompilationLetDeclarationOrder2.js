//// [tests/cases/compiler/jsFileCompilationLetDeclarationOrder2.ts] ////

//// [a.ts]
let b = 30;
a = 10;
//// [b.js]
let a = 10;
b = 30;


//// [b.js]
let a = 10;
b = 30;
//// [a.js]
let b = 30;
a = 10;
