//// [tests/cases/compiler/moduleResolutionWithSuffixes_one_jsonModule.ts] ////

//// [foo.ios.json]
{
	"ios": "platform ios"
}
//// [foo.json]
{
	"base": "platform base"
}

//// [index.ts]
import foo from "./foo.json";
console.log(foo.ios);

/bin/foo.ios.json(1,1): error TS1005: '{' expected.
/bin/foo.ios.json(1,2): error TS1136: Property assignment expected.
/bin/foo.ios.json(3,2): error TS1012: Unexpected token.
/bin/foo.ios.json(4,1): error TS1005: '}' expected.


==== /bin/foo.ios.json (4 errors) ====
    ({
    ~
!!! error TS1005: '{' expected.
     ~
!!! error TS1136: Property assignment expected.
        "ios": "platform ios"
    })
     ~
!!! error TS1012: Unexpected token.
    
    
!!! error TS1005: '}' expected.
//// [/bin/index.js]
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const foo_json_1 = __importDefault(require("./foo.json"));
console.log(foo_json_1.default.ios);
