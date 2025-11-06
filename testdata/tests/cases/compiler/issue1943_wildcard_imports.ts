// @module: commonjs
// @traceResolution: true

// Test with wildcard pattern in imports
// @Filename: /node_modules/pkg-test/package.json
{
  "name": "pkg-test",
  "type": "commonjs",
  "imports": {
    "#pkg-test/*.ts": "./*.d.ts"
  }
}

// @Filename: /node_modules/pkg-test/foo.d.ts
export declare const FOO: string;

// @Filename: /node_modules/pkg-test/bar.d.ts
export * from "#pkg-test/foo.ts";

// @Filename: /index.ts
import { FOO } from "pkg-test/bar";
console.log(FOO);