//// [tests/cases/conformance/node/nodePackageSelfName.ts] ////

//// [index.ts]
// esm format file
import * as self from "package";
self;
//// [index.mts]
// esm format file
import * as self from "package";
self;
//// [index.cts]
// esm format file
import * as self from "package";
self;
//// [package.json]
{
    "name": "package",
    "private": true,
    "type": "module",
    "exports": "./index.js"
}

//// [index.js]
// esm format file
import * as self from "package";
self;
//// [index.mjs]
// esm format file
import * as self from "package";
self;
//// [index.cjs]
// esm format file
import * as self from "package";
self;
