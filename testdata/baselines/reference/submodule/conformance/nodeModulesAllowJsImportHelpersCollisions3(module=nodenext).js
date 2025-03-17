//// [tests/cases/conformance/node/allowJs/nodeModulesAllowJsImportHelpersCollisions3.ts] ////

//// [index.js]
// cjs format file
export {default} from "fs";
export {default as foo} from "fs";
export {bar as baz} from "fs";
//// [index.js]
// esm format file
export {default} from "fs";
export {default as foo} from "fs";
export {bar as baz} from "fs";
//// [package.json]
{
    "name": "package",
    "private": true,
    "type": "module"
}
//// [package.json]
{
    "type": "commonjs"
}
//// [types.d.ts]
declare module "fs";
declare module "tslib" {
    export {};
    // intentionally missing all helpers
}

//// [index.js]
export { default } from "fs";
export { default as foo } from "fs";
export { bar as baz } from "fs";
//// [index.js]
export { default } from "fs";
export { default as foo } from "fs";
export { bar as baz } from "fs";
