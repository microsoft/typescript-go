//// [tests/cases/compiler/issue1943_export_star_reexport.ts] ////

//// [package.json]
{
  "name": "pkg-exporter-src",
  "type": "commonjs",
  "exports": {
    "./testing": {
      "types": "./dist/testing.d.ts",
      "default": "./src/testing.ts"
    }
  },
  "imports": {
    "#pkg-exporter-src/*.ts": {
      "types": "./src/*.ts",
      "default": "./src/*.ts"
    }
  }
}

//// [dep.ts]
export function stubDynamicConfig() {
  return "config";
}

//// [testing.ts]
// Re-export using import map pattern
export * from "#pkg-exporter-src/dep.ts";

//// [index.ts]
// This should work but may fail with: "Module has no exported member 'stubDynamicConfig'"
import { stubDynamicConfig } from "pkg-exporter-src/testing";

const result = stubDynamicConfig();
console.log(result);


//// [dep.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.stubDynamicConfig = stubDynamicConfig;
function stubDynamicConfig() {
    return "config";
}
//// [testing.js]
"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __exportStar = (this && this.__exportStar) || function(m, exports) {
    for (var p in m) if (p !== "default" && !Object.prototype.hasOwnProperty.call(exports, p)) __createBinding(exports, m, p);
};
Object.defineProperty(exports, "__esModule", { value: true });
// Re-export using import map pattern
__exportStar(require("#pkg-exporter-src/dep.ts"), exports);
//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// This should work but may fail with: "Module has no exported member 'stubDynamicConfig'"
const testing_1 = require("pkg-exporter-src/testing");
const result = (0, testing_1.stubDynamicConfig)();
console.log(result);


//// [dep.d.ts]
export declare function stubDynamicConfig(): string;
//// [testing.d.ts]
export * from "#pkg-exporter-src/dep.ts";
//// [index.d.ts]
export {};


//// [DtsFileErrors]


/node_modules/pkg-exporter-src/src/testing.d.ts(1,15): error TS2307: Cannot find module '#pkg-exporter-src/dep.ts' or its corresponding type declarations.


==== /node_modules/pkg-exporter-src/package.json (0 errors) ====
    {
      "name": "pkg-exporter-src",
      "type": "commonjs",
      "exports": {
        "./testing": {
          "types": "./dist/testing.d.ts",
          "default": "./src/testing.ts"
        }
      },
      "imports": {
        "#pkg-exporter-src/*.ts": {
          "types": "./src/*.ts",
          "default": "./src/*.ts"
        }
      }
    }
    
==== /node_modules/pkg-exporter-src/src/dep.d.ts (0 errors) ====
    export declare function stubDynamicConfig(): string;
    
==== /node_modules/pkg-exporter-src/src/testing.d.ts (1 errors) ====
    export * from "#pkg-exporter-src/dep.ts";
                  ~~~~~~~~~~~~~~~~~~~~~~~~~~
!!! error TS2307: Cannot find module '#pkg-exporter-src/dep.ts' or its corresponding type declarations.
    
==== /index.d.ts (0 errors) ====
    export {};
    