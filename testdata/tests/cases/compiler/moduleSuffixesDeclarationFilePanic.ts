// Test for panic when using moduleSuffixes with ".d" suffix
// @module: esnext
// @moduleResolution: bundler
// @moduleSuffixes: .d,

// @filename: /node_modules/my-lib/package.json
{
    "name": "my-lib",
    "version": "1.0.0",
    "main": "./index.js",
    "types": "./index.ts"
}

// @filename: /node_modules/my-lib/index.d.ts
export declare function bar(): void;

// @filename: /index.ts
import { bar } from "my-lib";
bar();
