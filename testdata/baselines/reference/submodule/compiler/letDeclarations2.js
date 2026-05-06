//// [tests/cases/compiler/letDeclarations2.ts] ////

//// [letDeclarations2.ts]
namespace M {
    let l1 = "s";
    export let l2 = 0;
}

//// [letDeclarations2.js]
"use strict";
var M;
(function (M) {
    let l1 = "s";
    M.l2 = 0;
})(M || (M = {}));


//// [letDeclarations2.d.ts]
namespace M {
    let l2: number;
}


//// [DtsFileErrors]


letDeclarations2.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== letDeclarations2.d.ts (1 errors) ====
    namespace M {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        let l2: number;
    }
    