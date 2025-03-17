//// [tests/cases/conformance/salsa/commonJSImportClassTypeReference.ts] ////

//// [main.js]
const { K } = require("./mod1");
/** @param {K} k */
function f(k) {
    k.values()
}

//// [mod1.js]
class K {
    values() {
        return new K()
    }
}
exports.K = K;


//// [main.js]
const { K } = require("./mod1");
function f(k) {
    k.values();
}
//// [mod1.js]
class K {
    values() {
        return new K();
    }
}
exports.K = K;
