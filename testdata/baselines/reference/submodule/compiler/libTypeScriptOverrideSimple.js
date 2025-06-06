//// [tests/cases/compiler/libTypeScriptOverrideSimple.ts] ////

//// [index.d.ts]
interface ABC { abc: string }
//// [index.ts]
/// <reference lib="dom" />
const a: ABC = { abc: "Hello" }

// This should fail because libdom has been replaced
// by the module above ^
window.localStorage


//// [index.js]
/// <reference lib="dom" />
const a = { abc: "Hello" };
// This should fail because libdom has been replaced
// by the module above ^
window.localStorage;
