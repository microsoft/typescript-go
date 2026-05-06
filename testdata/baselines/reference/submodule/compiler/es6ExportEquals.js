//// [tests/cases/compiler/es6ExportEquals.ts] ////

//// [es6ExportEquals.ts]
export function f() { }

export = f;


//// [es6ExportEquals.js]
export function f() { }


//// [es6ExportEquals.d.ts]
export function f(): void;
export = f;
