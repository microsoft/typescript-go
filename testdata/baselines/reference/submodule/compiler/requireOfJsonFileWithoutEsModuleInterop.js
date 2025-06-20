//// [tests/cases/compiler/requireOfJsonFileWithoutEsModuleInterop.ts] ////

//// [file1.ts]
import * as test from "./test.json"

//// [test.json]
{
    "a": true,
    "b": "hello"
}

out/test.json(1,1): error TS1005: '{' expected.
out/test.json(1,2): error TS1136: Property assignment expected.
out/test.json(4,2): error TS1012: Unexpected token.
out/test.json(5,1): error TS1005: '}' expected.


==== out/test.json (4 errors) ====
    ({
    ~
!!! error TS1005: '{' expected.
     ~
!!! error TS1136: Property assignment expected.
        "a": true,
        "b": "hello"
    })
     ~
!!! error TS1012: Unexpected token.
    
    
!!! error TS1005: '}' expected.
//// [out/file1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
