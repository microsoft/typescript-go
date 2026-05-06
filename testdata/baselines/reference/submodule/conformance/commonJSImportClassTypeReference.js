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


//// [mod1.js]
"use strict";
class K {
    values() {
        return new K();
    }
}
exports.K = K;
//// [main.js]
"use strict";
const { K } = require("./mod1");
/** @param {K} k */
function f(k) {
    k.values();
}


//// [mod1.d.ts]
class K {
    values(): K;
}
export { K };
//// [main.d.ts]
export {};


//// [DtsFileErrors]


out/mod1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/main.d.ts (0 errors) ====
    export {};
    
==== out/mod1.d.ts (1 errors) ====
    class K {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        values(): K;
    }
    export { K };
    