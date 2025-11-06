// @module: commonjs
// @declaration: true
// @traceResolution: true

// Setup pkg-exporter package with built declaration files
// @Filename: /node_modules/pkg-exporter/package.json
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

// @Filename: /node_modules/pkg-exporter/dist/dep.d.ts
export declare function stubDynamicConfig(): string;

// @Filename: /node_modules/pkg-exporter/dist/testing.d.ts
export * from "#pkg-exporter/dep.ts";

// @Filename: /index.ts
import { stubDynamicConfig } from "pkg-exporter/testing";

const result = stubDynamicConfig();
console.log(result);