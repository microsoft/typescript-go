//// [tests/cases/compiler/letDeclarations-useBeforeDefinition2.ts] ////

//// [file1.ts]
l;

//// [file2.ts]
const l = 0;


//// [file2.js]
const l = 0;
//// [file1.js]
l;
