// @module: commonjs
// @traceResolution: true

// Test with conditional wildcard pattern in imports
// @Filename: /node_modules/pkg-test/package.json
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

// @Filename: /node_modules/pkg-test/dist/foo.d.ts
export declare const FOO: string;

// @Filename: /node_modules/pkg-test/dist/bar.d.ts
export * from "#pkg-test/foo.ts";

// @Filename: /index.ts
import { FOO } from "pkg-test/dist/bar";
console.log(FOO);