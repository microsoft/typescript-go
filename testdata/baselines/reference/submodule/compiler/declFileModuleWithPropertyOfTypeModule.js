//// [tests/cases/compiler/declFileModuleWithPropertyOfTypeModule.ts] ////

//// [declFileModuleWithPropertyOfTypeModule.ts]
namespace m {
    export class c {
    }

    export var a = m;
}

//// [declFileModuleWithPropertyOfTypeModule.js]
"use strict";
var m;
(function (m) {
    class c {
    }
    m.c = c;
    m.a = m;
})(m || (m = {}));


//// [declFileModuleWithPropertyOfTypeModule.d.ts]
namespace m {
    class c {
    }
    var a: typeof m;
}


//// [DtsFileErrors]


declFileModuleWithPropertyOfTypeModule.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileModuleWithPropertyOfTypeModule.d.ts (1 errors) ====
    namespace m {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        class c {
        }
        var a: typeof m;
    }
    