// @module: commonjs
// @declaration: true

// Setup pkg-exporter-src package with exports and imports (source files)
// @Filename: /node_modules/pkg-exporter-src/package.json
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

// @Filename: /node_modules/pkg-exporter-src/src/dep.ts
export function stubDynamicConfig() {
  return "config";
}

// @Filename: /node_modules/pkg-exporter-src/src/testing.ts
// Re-export using import map pattern
export * from "#pkg-exporter-src/dep.ts";

// @Filename: /index.ts
// This should work but may fail with: "Module has no exported member 'stubDynamicConfig'"
import { stubDynamicConfig } from "pkg-exporter-src/testing";

const result = stubDynamicConfig();
console.log(result);
