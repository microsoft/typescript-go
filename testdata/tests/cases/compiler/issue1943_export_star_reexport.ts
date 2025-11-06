// @module: commonjs
// @declaration: true

// Setup pkg-exporter package with exports and imports
// @Filename: /node_modules/pkg-exporter/package.json
{
  "name": "pkg-exporter",
  "type": "commonjs",
  "exports": {
    "./testing": {
      "types": "./dist/testing.d.ts",
      "default": "./src/testing.ts"
    }
  },
  "imports": {
    "#pkg-exporter/*.ts": {
      "types": "./src/*.ts",
      "default": "./src/*.ts"
    }
  }
}

// @Filename: /node_modules/pkg-exporter/src/dep.ts
export function stubDynamicConfig() {
  return "config";
}

// @Filename: /node_modules/pkg-exporter/src/testing.ts
// Re-export using import map pattern
export * from "#pkg-exporter/dep.ts";

// @Filename: /index.ts
// This should work but may fail with: "Module has no exported member 'stubDynamicConfig'"
import { stubDynamicConfig } from "pkg-exporter/testing";

const result = stubDynamicConfig();
console.log(result);
