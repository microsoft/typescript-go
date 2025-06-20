//// [tests/cases/compiler/requireOfJsonFileWithComputedPropertyName.ts] ////

//// [file1.ts]
import b1 = require('./b.json');
let x = b1;
import b2 = require('./b.json');
if (x) {
    x = b2;
}

//// [b.json]
{
    [a]: 10
}

out/b.json(1,1): error TS1005: '{' expected.
out/b.json(1,2): error TS1136: Property assignment expected.
out/b.json(3,2): error TS1012: Unexpected token.
out/b.json(4,1): error TS1005: '}' expected.


==== out/b.json (4 errors) ====
    ({
    ~
!!! error TS1005: '{' expected.
     ~
!!! error TS1136: Property assignment expected.
        [a]: 10
    })
     ~
!!! error TS1012: Unexpected token.
    
    
!!! error TS1005: '}' expected.
//// [file1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const b1 = require("./b.json");
let x = b1;
const b2 = require("./b.json");
if (x) {
    x = b2;
}
