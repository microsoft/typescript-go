//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsExportDoubleAssignmentInClosure.ts] ////

//// [index.js]
// @ts-nocheck
function foo() {
    module.exports = exports = function (o) {
        return (o == null) ? create(base) : defineProperties(Object(o), descriptors);
    };
    const m = function () {
        // I have no idea what to put here
    }
    exports.methods = m;
}


//// [index.js]
function foo() {
    module.exports = exports = function (o) {
        return (o == null) ? create(base) : defineProperties(Object(o), descriptors);
    };
    const m = function () {
    };
    exports.methods = m;
}
