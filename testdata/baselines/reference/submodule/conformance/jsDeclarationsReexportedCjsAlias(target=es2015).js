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

//// [lib.js]
"use strict";
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
};
//// [main.js]
"use strict";
const { SomeClass, SomeClass: Another } = require('./lib');
module.exports = {
    SomeClass,
    Another
};


//// [lib.d.ts]
/**
 * @param {string} a
 */
function bar(a: string): string;
class SomeClass {
    a(): number;
}
const _default: {
    bar: typeof bar;
    SomeClass: typeof SomeClass;
};
export = _default;
//// [main.d.ts]
const _default: {
    SomeClass: {
        new (): {
            a(): number;
        };
    };
    Another: {
        new (): {
            a(): number;
        };
    };
};
export = _default;


//// [DtsFileErrors]


out/main.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/main.d.ts (1 errors) ====
    const _default: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        SomeClass: {
            new (): {
                a(): number;
            };
        };
        Another: {
            new (): {
                a(): number;
            };
        };
    };
    export = _default;
    
==== out/lib.d.ts (0 errors) ====
    /**
     * @param {string} a
     */
    function bar(a: string): string;
    class SomeClass {
        a(): number;
    }
    const _default: {
        bar: typeof bar;
        SomeClass: typeof SomeClass;
    };
    export = _default;
    