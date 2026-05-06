//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsTypeReferences2.ts] ////

//// [something.ts]
export const o = {
    a: 1,
    m: 1
}

//// [index.js]
const{ a, m } = require("./something").o;

const thing = a + m

module.exports = {
    thing
};


//// [something.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.o = void 0;
exports.o = {
    a: 1,
    m: 1
};
//// [index.js]
"use strict";
const { a, m } = require("./something").o;
const thing = a + m;
module.exports = {
    thing
};


//// [something.d.ts]
export const o: {
    a: number;
    m: number;
};
//// [index.d.ts]
const _default: {
    thing: number;
};
export = _default;


//// [DtsFileErrors]


tests/cases/conformance/jsdoc/declarations/out/index.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== tests/cases/conformance/jsdoc/declarations/out/index.d.ts (1 errors) ====
    const _default: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        thing: number;
    };
    export = _default;
    
==== tests/cases/conformance/jsdoc/declarations/out/something.d.ts (0 errors) ====
    export const o: {
        a: number;
        m: number;
    };
    