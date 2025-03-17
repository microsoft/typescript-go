//// [tests/cases/compiler/jsFileCompilationWithOutDeclarationFileNameSameAsInputJsFile.ts] ////

//// [a.ts]
class c {
}

//// [b.d.ts]
declare function foo(): boolean;


//// [a.js]
class c {
}
