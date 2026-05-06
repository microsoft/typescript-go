//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsCommonjsRelativePath.ts] ////

//// [thing.js]
'use strict';
class Thing {}
module.exports = { Thing }

//// [reexport.js]
'use strict';
const Thing = require('./thing').Thing
module.exports = { Thing }




//// [thing.d.ts]
class Thing {
}
const _default: {
    Thing: typeof Thing;
};
export = _default;
//// [reexport.d.ts]
const _default: {
    Thing: {
        new (): {};
    };
};
export = _default;


//// [DtsFileErrors]


reexport.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== reexport.d.ts (1 errors) ====
    const _default: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        Thing: {
            new (): {};
        };
    };
    export = _default;
    
==== thing.d.ts (0 errors) ====
    class Thing {
    }
    const _default: {
        Thing: typeof Thing;
    };
    export = _default;
    