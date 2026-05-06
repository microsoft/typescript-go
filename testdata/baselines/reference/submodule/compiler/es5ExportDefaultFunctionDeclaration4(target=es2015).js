//// [tests/cases/compiler/es5ExportDefaultFunctionDeclaration4.ts] ////

//// [es5ExportDefaultFunctionDeclaration4.ts]
declare module "bar" {
    var before: typeof func;

    export default function func(): typeof func;

    var after: typeof func;
}

//// [es5ExportDefaultFunctionDeclaration4.js]
"use strict";


//// [es5ExportDefaultFunctionDeclaration4.d.ts]
module "bar" {
    var before: typeof func;
    export default function func(): typeof func;
    var after: typeof func;
}


//// [DtsFileErrors]


es5ExportDefaultFunctionDeclaration4.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== es5ExportDefaultFunctionDeclaration4.d.ts (1 errors) ====
    module "bar" {
    ~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        var before: typeof func;
        export default function func(): typeof func;
        var after: typeof func;
    }
    