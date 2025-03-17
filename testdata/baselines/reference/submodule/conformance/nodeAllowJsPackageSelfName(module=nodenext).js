//// [tests/cases/conformance/node/allowJs/nodeAllowJsPackageSelfName.ts] ////

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
//// [package.json]
{
    "name": "package",
    "private": true,
    "type": "module",
    "exports": "./index.js"
}

//// [index.cjs]
import * as self from "package";
self;
//// [index.mjs]
import * as self from "package";
self;
//// [index.js]
import * as self from "package";
self;
