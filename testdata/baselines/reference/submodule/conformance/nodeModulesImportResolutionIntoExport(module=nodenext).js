//// [tests/cases/conformance/node/nodeModulesImportResolutionIntoExport.ts] ////

//// [index.ts]
// esm format file
import * as type from "#type";
type;
//// [index.mts]
// esm format file
import * as type from "#type";
type;
//// [index.cts]
// esm format file
import * as type from "#type";
type;
//// [package.json]
{
    "name": "package",
    "private": true,
    "type": "module",
    "exports": "./index.cjs",
    "imports": {
        "#type": "package"
    }
}

//// [index.cjs]
// esm format file
import * as type from "#type";
type;
//// [index.js]
// esm format file
import * as type from "#type";
type;
//// [index.mjs]
// esm format file
import * as type from "#type";
type;
