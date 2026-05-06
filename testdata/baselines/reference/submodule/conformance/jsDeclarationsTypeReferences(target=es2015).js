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
"use strict";
/// <reference types="node" />
const Something = require("fs").Something;
const thing = new Something();
module.exports = {
    thing
};


//// [index.d.ts]
const _default: {
    thing: import("fs").Something;
};
export = _default;


//// [DtsFileErrors]


tests/cases/conformance/jsdoc/declarations/out/index.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== tests/cases/conformance/jsdoc/declarations/out/index.d.ts (1 errors) ====
    const _default: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        thing: import("fs").Something;
    };
    export = _default;
    
==== node_modules/@types/node/index.d.ts (0 errors) ====
    declare module "fs" {
        export class Something {}
    }