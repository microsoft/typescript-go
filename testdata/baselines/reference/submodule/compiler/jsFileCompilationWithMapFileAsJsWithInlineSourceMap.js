//// [tests/cases/compiler/jsFileCompilationWithMapFileAsJsWithInlineSourceMap.ts] ////

//// [a.ts]
class c {
}

//// [b.js.map]
function foo() {
}

//// [b.js]
function bar() {
}

//// [b.js]
function bar() {
}
//// [a.js]
class c {
}
