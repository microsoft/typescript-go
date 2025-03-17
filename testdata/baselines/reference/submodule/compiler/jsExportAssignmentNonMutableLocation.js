//// [tests/cases/compiler/jsExportAssignmentNonMutableLocation.ts] ////

//// [file.js]
const customSymbol = Symbol("custom");

// This is a common pattern in Node’s built-in modules:
module.exports = {
    customSymbol,
};

exports.customSymbol2 = Symbol("custom");

//// [file.js]
const customSymbol = Symbol("custom");
module.exports = {
    customSymbol,
};
exports.customSymbol2 = Symbol("custom");
