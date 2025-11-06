// @module: commonjs
// @declaration: true

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

// Built declaration files (not source)
// @Filename: /node_modules/pkg-exporter/dist/dep.d.ts
export declare function stubDynamicConfig(): string;

// @Filename: /node_modules/pkg-exporter/dist/testing.d.ts
// Re-export using import map pattern - but now in a .d.ts file
export * from "#pkg-exporter/dep.ts";

// @Filename: /index.ts
// This should work but may fail with: "Module has no exported member 'stubDynamicConfig'"
import { stubDynamicConfig } from "pkg-exporter/testing";

const result = stubDynamicConfig();
console.log(result);