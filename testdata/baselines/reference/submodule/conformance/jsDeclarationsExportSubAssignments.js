//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsExportSubAssignments.ts] ////

//// [cls.js]
const Strings = {
    a: "A",
    b: "B"
};
class Foo {}
module.exports = Foo;
module.exports.Strings = Strings;

//// [cls.js]
const Strings = {
    a: "A",
    b: "B"
};
class Foo {
}
module.exports = Foo;
module.exports.Strings = Strings;
