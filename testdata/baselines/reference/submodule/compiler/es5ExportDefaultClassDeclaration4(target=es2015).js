//// [tests/cases/compiler/es5ExportDefaultClassDeclaration4.ts] ////

//// [es5ExportDefaultClassDeclaration4.ts]
declare module "foo" {
    export var before: C;

    export default class C {
        method(): C;
    }

    export var after: C;

    export var t: typeof C;
}



//// [es5ExportDefaultClassDeclaration4.js]
"use strict";


//// [es5ExportDefaultClassDeclaration4.d.ts]
module "foo" {
    var before: C;
    export default class C {
        method(): C;
    }
    var after: C;
    var t: typeof C;
}


//// [DtsFileErrors]


es5ExportDefaultClassDeclaration4.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== es5ExportDefaultClassDeclaration4.d.ts (1 errors) ====
    module "foo" {
    ~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        var before: C;
        export default class C {
            method(): C;
        }
        var after: C;
        var t: typeof C;
    }
    