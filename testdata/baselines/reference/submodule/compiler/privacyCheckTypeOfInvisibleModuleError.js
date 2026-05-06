//// [tests/cases/compiler/privacyCheckTypeOfInvisibleModuleError.ts] ////

//// [privacyCheckTypeOfInvisibleModuleError.ts]
namespace Outer {
    namespace Inner {
        export var m: typeof Inner;
    }

    export var f: typeof Inner;
}


//// [privacyCheckTypeOfInvisibleModuleError.js]
"use strict";
var Outer;
(function (Outer) {
    let Inner;
    (function (Inner) {
    })(Inner || (Inner = {}));
})(Outer || (Outer = {}));


//// [privacyCheckTypeOfInvisibleModuleError.d.ts]
namespace Outer {
    namespace Inner {
        var m: typeof Inner;
    }
    export var f: typeof Inner;
    export {};
}


//// [DtsFileErrors]


privacyCheckTypeOfInvisibleModuleError.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== privacyCheckTypeOfInvisibleModuleError.d.ts (1 errors) ====
    namespace Outer {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        namespace Inner {
            var m: typeof Inner;
        }
        export var f: typeof Inner;
        export {};
    }
    