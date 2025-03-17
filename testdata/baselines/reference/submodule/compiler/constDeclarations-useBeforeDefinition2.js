//// [tests/cases/compiler/constDeclarations-useBeforeDefinition2.ts] ////

//// [file1.ts]
c;

//// [file2.ts]
const c = 0;


//// [file2.js]
const c = 0;
//// [file1.js]
c;
