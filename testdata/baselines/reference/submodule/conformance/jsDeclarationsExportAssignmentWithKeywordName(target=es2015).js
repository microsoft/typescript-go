//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsExportAssignmentWithKeywordName.ts] ////

//// [index.js]
var x = 12;
module.exports = {
    extends: 'base',
    more: {
        others: ['strs']
    },
    x
};

//// [index.js]
"use strict";
var x = 12;
module.exports = {
    extends: 'base',
    more: {
        others: ['strs']
    },
    x
};


//// [index.d.ts]
const _default: {
    extends: string;
    more: {
        others: string[];
    };
    x: number;
};
export = _default;


//// [DtsFileErrors]


out/index.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/index.d.ts (1 errors) ====
    const _default: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        extends: string;
        more: {
            others: string[];
        };
        x: number;
    };
    export = _default;
    