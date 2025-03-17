//// [tests/cases/conformance/node/allowJs/nodeModulesAllowJsExportAssignment.ts] ////

//// [index.js]
// cjs format file
const a = {};
export = a;
//// [file.js]
// cjs format file
const a = {};
module.exports = a;
//// [index.js]
// esm format file
const a = {};
export = a;
//// [file.js]
// esm format file
import "fs";
const a = {};
module.exports = a;
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

//// [index.js]
const a = {};
export {};
//// [file.js]
const a = {};
module.exports = a;
//// [index.js]
const a = {};
export {};
//// [file.js]
import "fs";
const a = {};
module.exports = a;
