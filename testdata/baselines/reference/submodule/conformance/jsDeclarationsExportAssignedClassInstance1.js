//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsExportAssignedClassInstance1.ts] ////

//// [index.js]
class Foo {}

module.exports = new Foo();

//// [index.js]
class Foo {
}
module.exports = new Foo();
