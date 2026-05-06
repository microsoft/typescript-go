//// [tests/cases/conformance/es6/yieldExpressions/generatorInAmbientContext6.ts] ////

//// [generatorInAmbientContext6.ts]
namespace M {
    export function *generator(): any { }
}

//// [generatorInAmbientContext6.js]
"use strict";
var M;
(function (M) {
    function* generator() { }
    M.generator = generator;
})(M || (M = {}));


//// [generatorInAmbientContext6.d.ts]
namespace M {
    function generator(): any;
}


//// [DtsFileErrors]


generatorInAmbientContext6.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== generatorInAmbientContext6.d.ts (1 errors) ====
    namespace M {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        function generator(): any;
    }
    