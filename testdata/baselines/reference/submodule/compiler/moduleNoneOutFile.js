//// [tests/cases/compiler/moduleNoneOutFile.ts] ////

//// [first.ts]
class Foo {}
//// [second.ts]
class Bar extends Foo {}

//// [first.js]
class Foo {
}
//// [second.js]
class Bar extends Foo {
}
