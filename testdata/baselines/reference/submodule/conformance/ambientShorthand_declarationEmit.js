//// [tests/cases/conformance/ambient/ambientShorthand_declarationEmit.ts] ////

//// [ambientShorthand_declarationEmit.ts]
declare module "foo";


//// [ambientShorthand_declarationEmit.js]
"use strict";


//// [ambientShorthand_declarationEmit.d.ts]
module "foo";


//// [DtsFileErrors]


ambientShorthand_declarationEmit.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== ambientShorthand_declarationEmit.d.ts (1 errors) ====
    module "foo";
    ~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    