//// [tests/cases/compiler/declarationEmitHasTypesRefOnNamespaceUse.ts] ////

//// [dep.d.ts]
declare namespace NS {
    interface Dep {
    }
}
//// [package.json]
{
    "typings": "dep.d.ts"
}
//// [index.ts]
class Src implements NS.Dep { }


//// [index.js]
"use strict";
class Src {
}


//// [index.d.ts]
class Src implements NS.Dep {
}


//// [DtsFileErrors]


/src/index.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== /src/index.d.ts (1 errors) ====
    class Src implements NS.Dep {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    
==== /deps/dep/dep.d.ts (0 errors) ====
    declare namespace NS {
        interface Dep {
        }
    }
==== /deps/dep/package.json (0 errors) ====
    {
        "typings": "dep.d.ts"
    }