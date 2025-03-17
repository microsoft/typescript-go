//// [tests/cases/compiler/declarationFileOverwriteErrorWithOut.ts] ////

//// [out.d.ts]
declare class c {
}

//// [a.ts]
class d {
}


//// [a.js]
class d {
}
