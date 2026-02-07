currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/common/src/greeter.ts] *new* 
export function greet(name: string): string {
    return "Hello, " + name + "!";
}
//// [/home/src/workspaces/project/sub/src/index.ts] *new* 
import { greet } from "../../common/src/greeter.js";
console.log(greet("world"));
//// [/home/src/workspaces/project/sub/tsconfig.json] *new* 
{
    "extends": "../tsconfig.json",
    "include": [
        "src/**/*",
        "../common/src/**/*"
    ]
}
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{
    "compilerOptions": {
        "target": "ES2024",
        "module": "NodeNext",
        "moduleResolution": "NodeNext",
        "outDir": "${configDir}/dist/",
        "strict": true,
        "skipLibCheck": true
    },
    "include": ["common/src/**/*"]
}

tsgo -p sub/tsconfig.json --explainFiles
ExitStatus:: Success
Output::
../../tslibs/TS/Lib/lib.es2024.full.d.ts
   Default library for target 'ES2024'
common/src/greeter.ts
   Imported via "../../common/src/greeter.js" from file 'sub/src/index.ts'
   Matched by include pattern '../common/src/**/*' in 'sub/tsconfig.json'
   File is CommonJS module because 'package.json' was not found
sub/src/index.ts
   Matched by include pattern 'src/**/*' in 'sub/tsconfig.json'
   File is CommonJS module because 'package.json' was not found
//// [/home/src/tslibs/TS/Lib/lib.es2024.full.d.ts] *Lib*
/// <reference no-default-lib="true"/>
interface Boolean {}
interface Function {}
interface CallableFunction {}
interface NewableFunction {}
interface IArguments {}
interface Number { toExponential: any; }
interface Object {}
interface RegExp {}
interface String { charAt: any; }
interface Array<T> { length: number; [n: number]: T; }
interface ReadonlyArray<T> {}
interface SymbolConstructor {
    (desc?: string | number): symbol;
    for(name: string): symbol;
    readonly toStringTag: symbol;
}
declare var Symbol: SymbolConstructor;
interface Symbol {
    readonly [Symbol.toStringTag]: string;
}
declare const console: { log(msg: any): void; };
//// [/home/src/workspaces/project/common/src/greeter.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.greet = greet;
function greet(name) {
    return "Hello, " + name + "!";
}

//// [/home/src/workspaces/project/sub/dist/src/index.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const greeter_js_1 = require("../../common/src/greeter.js");
console.log((0, greeter_js_1.greet)("world"));


