//// [tests/cases/conformance/salsa/commonJSImportNestedClassTypeReference.ts] ////

//// [main.js]
const { K } = require("./mod1");
/** @param {K} k */
function f(k) {
    k.values()
}

//// [mod1.js]
var NS = {}
NS.K =class {
    values() {
        return new NS.K()
    }
}
exports.K = NS.K;


//// [main.js]
const { K } = require("./mod1");
function f(k) {
    k.values();
}
//// [mod1.js]
var NS = {};
NS.K = class {
    values() {
        return new NS.K();
    }
};
exports.K = NS.K;
