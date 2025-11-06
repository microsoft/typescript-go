//// [tests/cases/compiler/issue1943_conditional_imports.ts] ////

//// [package.json]
{
  "name": "pkg-test",
  "type": "commonjs",
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
import { FOO } from "pkg-test/dist/bar";
console.log(FOO);

//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const bar_1 = require("pkg-test/dist/bar");
console.log(bar_1.FOO);
