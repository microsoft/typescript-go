//// [tests/cases/compiler/jsExportAssignmentNonMutableLocation.ts] ////

//// [file.js]
const customSymbol = Symbol("custom");

// This is a common pattern in Node’s built-in modules:
module.exports = {
    customSymbol,
};

exports.customSymbol2 = Symbol("custom");



//// [file.d.ts]
const _default: {
    customSymbol: symbol;
};
export = _default;
export var customSymbol2: symbol;
