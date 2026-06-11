//// [tests/cases/compiler/cjsExportsDuplication.ts] ////

//// [file.js]
exports.foo = 42
exports.foo = "hello"
exports.foo = true



//// [file.d.ts]
export declare var foo: "hello" | 42 | true;
