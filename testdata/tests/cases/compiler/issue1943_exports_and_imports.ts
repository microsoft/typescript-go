// @module: commonjs
// @declaration: true
// @traceResolution: true

// Test with exports AND imports together
// @Filename: /node_modules/pkg-test/package.json
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

// @Filename: /node_modules/pkg-test/dist/foo.d.ts
export declare const FOO: string;

// @Filename: /node_modules/pkg-test/dist/bar.d.ts
export * from "#pkg-test/foo.ts";

// @Filename: /index.ts
// Import through exports path
import { FOO } from "pkg-test/bar";
console.log(FOO);