//// [tests/cases/compiler/privacyCheckTypeOfInvisibleModuleNoError.ts] ////

//// [privacyCheckTypeOfInvisibleModuleNoError.ts]
namespace Outer {
    namespace Inner {
        export var m: number;
    }

    export var f: typeof Inner; // Since we dont unwind inner any more, it is error here
}


//// [privacyCheckTypeOfInvisibleModuleNoError.js]
"use strict";
var Outer;
(function (Outer) {
    let Inner;
    (function (Inner) {
    })(Inner || (Inner = {}));
})(Outer || (Outer = {}));


//// [privacyCheckTypeOfInvisibleModuleNoError.d.ts]
namespace Outer {
    namespace Inner {
        var m: number;
    }
    export var f: typeof Inner;
    export {};
}


//// [DtsFileErrors]


privacyCheckTypeOfInvisibleModuleNoError.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== privacyCheckTypeOfInvisibleModuleNoError.d.ts (1 errors) ====
    namespace Outer {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        namespace Inner {
            var m: number;
        }
        export var f: typeof Inner;
        export {};
    }
    