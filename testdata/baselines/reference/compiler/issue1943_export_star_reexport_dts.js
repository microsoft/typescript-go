//// [tests/cases/compiler/issue1943_export_star_reexport_dts.ts] ////

//// [package.json]
{
  "name": "pkg-exporter",
  "type": "commonjs",
  "exports": {
    "./testing": {
      "types": "./dist/testing.d.ts",
      "default": "./dist/testing.js"
    }
  },
  "imports": {
    "#pkg-exporter/*.ts": {
      "types": "./dist/*.d.ts",
      "default": "./dist/*.js"
    }
  }
}

// Built declaration files (not source)
//// [dep.d.ts]
export declare function stubDynamicConfig(): string;

//// [testing.d.ts]
// Re-export using import map pattern - but now in a .d.ts file
export * from "#pkg-exporter/dep.ts";

//// [index.ts]
// This should work but may fail with: "Module has no exported member 'stubDynamicConfig'"
import { stubDynamicConfig } from "pkg-exporter/testing";

const result = stubDynamicConfig();
console.log(result);

//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// This should work but may fail with: "Module has no exported member 'stubDynamicConfig'"
const testing_1 = require("pkg-exporter/testing");
const result = (0, testing_1.stubDynamicConfig)();
console.log(result);


//// [index.d.ts]
export {};
