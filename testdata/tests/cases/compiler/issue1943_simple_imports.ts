// @module: commonjs
// @traceResolution: true

// Simpler test: Just test if imports pattern works at all
// @Filename: /node_modules/pkg-test/package.json
{
  "name": "pkg-test",
  "type": "commonjs",
  "imports": {
    "#pkg-test/foo.ts": "./foo.d.ts"
  }
}

// @Filename: /node_modules/pkg-test/foo.d.ts
export declare const FOO: string;

// @Filename: /node_modules/pkg-test/bar.d.ts
export * from "#pkg-test/foo.ts";

// @Filename: /index.ts
import { FOO } from "pkg-test/bar";
console.log(FOO);