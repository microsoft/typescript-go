//// [tests/cases/compiler/issue1943_exports_and_imports.ts] ////

//// [package.json]
{
  "name": "pkg-test",
  "type": "commonjs",
  "exports": {
    "./bar": {
      "types": "./dist/bar.d.ts",
      "default": "./dist/bar.js"
    }
  },
  "imports": {
    "#pkg-test/*.ts": {
      "types": "./dist/*.d.ts",
      "default": "./dist/*.js"
    }
  }
}

//// [foo.d.ts]
export declare const FOO: string;

//// [bar.d.ts]
export * from "#pkg-test/foo.ts";

//// [index.ts]
// Import through exports path
import { FOO } from "pkg-test/bar";
console.log(FOO);

//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// Import through exports path
const bar_1 = require("pkg-test/bar");
console.log(bar_1.FOO);


//// [index.d.ts]
export {};
