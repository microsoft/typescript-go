//// [tests/cases/compiler/isolatedModules_resolveJsonModule_strict_outDir_commonJs.ts] ////

//// [a.ts]
import * as j from "./j.json";

//// [j.json]
{}


/dist/j.json(1,1): error TS1005: '{' expected.
/dist/j.json(1,2): error TS1136: Property assignment expected.
/dist/j.json(1,4): error TS1012: Unexpected token.
/dist/j.json(2,1): error TS1005: '}' expected.


==== /dist/j.json (4 errors) ====
    ({})
    ~
!!! error TS1005: '{' expected.
     ~
!!! error TS1136: Property assignment expected.
       ~
!!! error TS1012: Unexpected token.
    
    
!!! error TS1005: '}' expected.
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
