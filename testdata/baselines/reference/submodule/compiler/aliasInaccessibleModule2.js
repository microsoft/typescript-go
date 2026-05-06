//// [tests/cases/compiler/aliasInaccessibleModule2.ts] ////

//// [aliasInaccessibleModule2.ts]
namespace M {
    namespace N {
        class C {
        }
        
    }
    import R = N;
    export import X = R;
}

//// [aliasInaccessibleModule2.js]
"use strict";
var M;
(function (M) {
    let N;
    (function (N) {
        class C {
        }
    })(N || (N = {}));
    var R = N;
    M.X = R;
})(M || (M = {}));


//// [aliasInaccessibleModule2.d.ts]
namespace M {
    namespace N {
    }
    import R = N;
    export import X = R;
    export {};
}


//// [DtsFileErrors]


aliasInaccessibleModule2.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== aliasInaccessibleModule2.d.ts (1 errors) ====
    namespace M {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        namespace N {
        }
        import R = N;
        export import X = R;
        export {};
    }
    