//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsTypeReassignmentFromDeclaration.ts] ////

//// [some-mod.d.ts]
interface Item {
    x: string;
}
declare const items: Item[];
export = items;
//// [index.js]
/** @type {typeof import("/some-mod")} */
const items = [];
module.exports = items;

//// [index.js]
"use strict";
/** @type {typeof import("/some-mod")} */
const items = [];
module.exports = items;


//// [index.d.ts]
/** @type {typeof import("/some-mod")} */
const items: typeof import("/some-mod");
export = items;


//// [DtsFileErrors]


/out/index.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== /some-mod.d.ts (0 errors) ====
    interface Item {
        x: string;
    }
    declare const items: Item[];
    export = items;
==== /out/index.d.ts (1 errors) ====
    /** @type {typeof import("/some-mod")} */
    const items: typeof import("/some-mod");
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export = items;
    