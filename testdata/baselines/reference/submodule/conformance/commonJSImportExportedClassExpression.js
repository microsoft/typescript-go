//// [tests/cases/conformance/salsa/commonJSImportExportedClassExpression.ts] ////

//// [main.js]
const { K } = require("./mod1");
/** @param {K} k */
function f(k) {
    k.values()
}

//// [mod1.js]
exports.K = class K {
    values() {
    }
};


//// [main.js]
const { K } = require("./mod1");
function f(k) {
    k.values();
}
//// [mod1.js]
exports.K = class K {
    values() {
    }
};
