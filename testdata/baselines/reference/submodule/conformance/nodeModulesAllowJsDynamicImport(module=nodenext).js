//// [tests/cases/conformance/node/allowJs/nodeModulesAllowJsDynamicImport.ts] ////

//// [index.js]
// cjs format file
export async function main() {
    const { readFile } = await import("fs");
}
//// [index.js]
// esm format file
export async function main() {
    const { readFile } = await import("fs");
}
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

//// [index.js]
export async function main() {
    const { readFile } = await import("fs");
}
//// [index.js]
export async function main() {
    const { readFile } = await import("fs");
}
