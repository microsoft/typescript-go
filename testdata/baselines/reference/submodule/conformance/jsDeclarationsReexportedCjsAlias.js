//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsReexportedCjsAlias.ts] ////

//// [lib.js]
/**
 * @param {string} a
 */
function bar(a) {
    return a + a;
}

class SomeClass {
    a() {
        return 1;
    }
}

module.exports = {
    bar,
    SomeClass
}
//// [main.js]
const { SomeClass, SomeClass: Another } = require('./lib');

module.exports = {
    SomeClass,
    Another
}

//// [main.js]
const { SomeClass, SomeClass: Another } = require('./lib');
module.exports = {
    SomeClass,
    Another
};
