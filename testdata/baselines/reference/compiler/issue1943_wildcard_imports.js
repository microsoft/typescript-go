//// [tests/cases/compiler/issue1943_wildcard_imports.ts] ////

//// [package.json]
{
  "name": "pkg-test",
  "type": "commonjs",
  "imports": {
    "#pkg-test/*.ts": "./*.d.ts"
  }
}

//// [foo.d.ts]
export declare const FOO: string;

//// [bar.d.ts]
export * from "#pkg-test/foo.ts";

//// [index.ts]
import { FOO } from "pkg-test/bar";
console.log(FOO);

//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const bar_1 = require("pkg-test/bar");
console.log(bar_1.FOO);
