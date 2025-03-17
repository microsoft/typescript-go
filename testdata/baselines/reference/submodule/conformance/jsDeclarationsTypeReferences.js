//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsTypeReferences.ts] ////

//// [index.d.ts]
declare module "fs" {
    export class Something {}
}
//// [index.js]
/// <reference types="node" />

const Something = require("fs").Something;

const thing = new Something();

module.exports = {
    thing
};


//// [index.js]
const Something = require("fs").Something;
const thing = new Something();
module.exports = {
    thing
};
