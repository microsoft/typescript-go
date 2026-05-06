//// [tests/cases/compiler/declarationEmitDestructuringPrivacyError.ts] ////

//// [declarationEmitDestructuringPrivacyError.ts]
namespace m {
    class c {
    }
    export var [x, y, z] = [10, new c(), 30];
}

//// [declarationEmitDestructuringPrivacyError.js]
"use strict";
var m;
(function (m) {
    var _a;
    class c {
    }
    _a = [10, new c(), 30], m.x = _a[0], m.y = _a[1], m.z = _a[2];
})(m || (m = {}));


//// [declarationEmitDestructuringPrivacyError.d.ts]
namespace m {
    class c {
    }
    export var x: number, y: c, z: number;
    export {};
}


//// [DtsFileErrors]


declarationEmitDestructuringPrivacyError.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitDestructuringPrivacyError.d.ts (1 errors) ====
    namespace m {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        class c {
        }
        export var x: number, y: c, z: number;
        export {};
    }
    