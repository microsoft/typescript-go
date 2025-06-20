//// [tests/cases/compiler/jsonFileImportChecksCallCorrectlyTwice.ts] ////

//// [index.ts]
import data from "./data.json";

interface Foo {
  str: string;
}

fn(data.foo);
fn(data.foo); // <-- shouldn't error!

function fn(arg: Foo[]) { }
//// [data.json]
{
    "foo": [
      {
        "bool": true,
        "str": "123"
      }
    ]
}

dist/data.json(1,1): error TS1005: '{' expected.
dist/data.json(1,2): error TS1136: Property assignment expected.
dist/data.json(8,2): error TS1012: Unexpected token.
dist/data.json(9,1): error TS1005: '}' expected.


==== dist/data.json (4 errors) ====
    ({
    ~
!!! error TS1005: '{' expected.
     ~
!!! error TS1136: Property assignment expected.
        "foo": [
            {
                "bool": true,
                "str": "123"
            }
        ]
    })
     ~
!!! error TS1012: Unexpected token.
    
    
!!! error TS1005: '}' expected.
//// [index.js]
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const data_json_1 = __importDefault(require("./data.json"));
fn(data_json_1.default.foo);
fn(data_json_1.default.foo); // <-- shouldn't error!
function fn(arg) { }
