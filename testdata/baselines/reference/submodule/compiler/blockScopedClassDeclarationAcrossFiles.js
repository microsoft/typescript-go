//// [tests/cases/compiler/blockScopedClassDeclarationAcrossFiles.ts] ////

//// [c.ts]
let foo: typeof C;
//// [b.ts]
class C { }


//// [b.js]
class C {
}
//// [c.js]
let foo;
