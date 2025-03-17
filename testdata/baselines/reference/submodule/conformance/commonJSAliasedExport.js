//// [tests/cases/conformance/salsa/commonJSAliasedExport.ts] ////

//// [commonJSAliasedExport.js]
const donkey = (ast) =>  ast;

function funky(declaration) {
    return false;
}
module.exports = donkey;
module.exports.funky = funky;

//// [bug43713.js]
const { funky } = require('./commonJSAliasedExport');
/** @type {boolean} */
var diddy
var diddy = funky(1)



//// [bug43713.js]
const { funky } = require('./commonJSAliasedExport');
var diddy;
var diddy = funky(1);
