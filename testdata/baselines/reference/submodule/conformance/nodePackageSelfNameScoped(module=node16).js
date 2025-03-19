//// [tests/cases/conformance/node/nodePackageSelfNameScoped.ts] ////

//// [index.ts]
// esm format file
import * as self from "@scope/package";
self;
//// [index.mts]
// esm format file
import * as self from "@scope/package";
self;
//// [index.cts]
// cjs format file
import * as self from "@scope/package";
self;
//// [package.json]
{
    "name": "@scope/package",
    "private": true,
    "type": "module",
    "exports": "./index.js"
}

//// [index.js]
// esm format file
import * as self from "@scope/package";
self;
//// [index.mjs]
// esm format file
import * as self from "@scope/package";
self;
//// [index.cjs]
// cjs format file
import * as self from "@scope/package";
self;
