//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsExportedClassAliases.ts] ////

//// [errors.js]
class FancyError extends Error {
    constructor(status) {
        super(`error with status ${status}`);
    }
}

module.exports = {
    FancyError
};

//// [index.js]
// issue arises here on compilation
const errors = require("./errors");

module.exports = {
    errors
};

//// [errors.js]
"use strict";
class FancyError extends Error {
    constructor(status) {
        super(`error with status ${status}`);
    }
}
module.exports = {
    FancyError
};
//// [index.js]
"use strict";
// issue arises here on compilation
const errors = require("./errors");
module.exports = {
    errors
};


//// [errors.d.ts]
class FancyError extends Error {
    constructor(status: any);
}
const _default: {
    FancyError: typeof FancyError;
};
export = _default;
//// [index.d.ts]
const _default: {
    errors: {
        FancyError: {
            new (status: any): {
                name: string;
                message: string;
                stack?: string;
            };
        };
    };
};
export = _default;


//// [DtsFileErrors]


out/index.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/index.d.ts (1 errors) ====
    const _default: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        errors: {
            FancyError: {
                new (status: any): {
                    name: string;
                    message: string;
                    stack?: string;
                };
            };
        };
    };
    export = _default;
    
==== out/errors.d.ts (0 errors) ====
    class FancyError extends Error {
        constructor(status: any);
    }
    const _default: {
        FancyError: typeof FancyError;
    };
    export = _default;
    