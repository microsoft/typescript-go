//// [tests/cases/compiler/pnpSimpleTest.ts] ////

//// [.pnp.cjs]
module.exports = {};

//// [.pnp.data.json]
{
  "dependencyTreeRoots": [
    {
      "name": "project",
      "reference": "workspace:."
    }
  ],
  "ignorePatternData": null,
  "enableTopLevelFallback": false,
  "fallbackPool": [],
  "fallbackExclusionList": [],
  "packageRegistryData": [
    ["project", [
      ["workspace:.", {
        "packageLocation": "./",
        "packageDependencies": [
          ["package-a", "npm:1.0.0"],
          ["package-b", "npm:2.0.0"]
        ]
      }]
    ]],
    ["package-a", [
      ["npm:1.0.0", {
        "packageLocation": "./.yarn/cache/package-a-npm-1.0.0-abcd1234/node_modules/package-a/",
        "packageDependencies": []
      }]
    ]],
    ["package-b", [
      ["npm:2.0.0", {
        "packageLocation": "./.yarn/cache/package-b-npm-2.0.0-efgh5678/node_modules/package-b/",
        "packageDependencies": []
      }]
    ]]
  ]
}

//// [package.json]
{
  "name": "project",
  "dependencies": {
    "package-a": "npm:1.0.0",
    "package-b": "npm:2.0.0"
  }
}

//// [package.json]
{
  "name": "package-a",
  "version": "1.0.0",
  "exports": {
    ".": "./index.js"
  },
  "types": "index.d.ts"
}

//// [index.d.ts]
export declare function helperA(value: string): string;

//// [package.json]
{
  "name": "package-b", 
  "version": "2.0.0",
  "exports": {
    ".": "./index.js"
  },
  "types": "index.d.ts"
}

//// [index.d.ts]
export declare function helperB(value: number): string;

//// [index.ts]
// Workspace package that imports both third-party dependencies
import { helperA } from 'package-a';
import { helperB } from 'package-b';

export function processData(text: string, num: number): string {
  const resultA = helperA(text);
  const resultB = helperB(num);
  return `${resultA} | ${resultB}`;
}

//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.processData = processData;
// Workspace package that imports both third-party dependencies
const package_a_1 = require("package-a");
const package_b_1 = require("package-b");
function processData(text, num) {
    const resultA = (0, package_a_1.helperA)(text);
    const resultB = (0, package_b_1.helperB)(num);
    return `${resultA} | ${resultB}`;
}


//// [index.d.ts]
export declare function processData(text: string, num: number): string;
