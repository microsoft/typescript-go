//// [tests/cases/compiler/aliasInaccessibleModule.ts] ////

//// [aliasInaccessibleModule.ts]
namespace M {
    namespace N {
    }
    export import X = N;
}

//// [aliasInaccessibleModule.js]
"use strict";
var M;
(function (M) {
})(M || (M = {}));


//// [aliasInaccessibleModule.d.ts]
namespace M {
    namespace N {
    }
    export import X = N;
    export {};
}


//// [DtsFileErrors]


aliasInaccessibleModule.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== aliasInaccessibleModule.d.ts (1 errors) ====
    namespace M {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        namespace N {
        }
        export import X = N;
        export {};
    }
    