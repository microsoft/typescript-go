
currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
//// [/home/src/workspaces/project/node_modules/@types/react/index.d.ts] *new* 
export {};
declare global {
    namespace JSX {
        interface Element {}
        interface IntrinsicElements {
            div: {
                propA?: boolean;
            };
        }
    }
}
//// [/home/src/workspaces/project/node_modules/react/jsx-runtime.js] *new* 
export {}
//// [/home/src/workspaces/project/src/index.tsx] *new* 
export const App = () => <div propA={true}></div>;
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{ 
    "compilerOptions": { 
        "module": "commonjs",
        "jsx": "react-jsx", 
        "incremental": true, 
        "jsxImportSource": "react" 
    } 
}

ExitStatus:: 0

CompilerOptions::{}
Output::
//// [/home/src/tslibs/TS/Lib/lib.d.ts] *Lib*
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
//// [/home/src/workspaces/project/src/index.js] *new* 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.App = void 0;
const jsx_runtime_1 = require("react/jsx-runtime");
const App = () => jsx_runtime_1.jsx("div", { propA: true });
exports.App = App;

//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./src/index.tsx","./node_modules/@types/react/index.d.ts"],"fileInfos":[{"version":"7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e","affectsGlobalScope":true,"impliedNodeFormat":1},"0f825b9ea7708adb1fe25efed7859af428d3ceb0d5df7f67f0a8c9eed854825a",{"version":"3301e1e26bf8906b220ec647394a80968ea8de76669750be485eeb41cf9e8679","affectsGlobalScope":true,"impliedNodeFormat":1}],"options":{"jsx":4,"jsxImportSource":"react","module":1}}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./src/index.tsx",
    "./node_modules/@types/react/index.d.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../tslibs/TS/Lib/lib.d.ts",
      "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e",
      "signature": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7dee939514de4bde7a51760a39e2b3bfa068bfc4a2939e1dbad2bfdf2dc4662e",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./src/index.tsx",
      "version": "0f825b9ea7708adb1fe25efed7859af428d3ceb0d5df7f67f0a8c9eed854825a",
      "signature": "0f825b9ea7708adb1fe25efed7859af428d3ceb0d5df7f67f0a8c9eed854825a",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./node_modules/@types/react/index.d.ts",
      "version": "3301e1e26bf8906b220ec647394a80968ea8de76669750be485eeb41cf9e8679",
      "signature": "3301e1e26bf8906b220ec647394a80968ea8de76669750be485eeb41cf9e8679",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "3301e1e26bf8906b220ec647394a80968ea8de76669750be485eeb41cf9e8679",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    }
  ],
  "options": {
    "jsx": 4,
    "jsxImportSource": "react",
    "module": 1
  },
  "size": 523
}


SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/src/index.tsx
*refresh*    /home/src/workspaces/project/node_modules/@types/react/index.d.ts

Signatures::
