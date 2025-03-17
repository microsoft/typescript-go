//// [tests/cases/conformance/node/allowJs/nodeModulesAllowJsTopLevelAwait.ts] ////

//// [index.js]
// cjs format file
const x = await 1;
export {x};
for await (const y of []) {}
//// [index.js]
// esm format file
const x = await 1;
export {x};
for await (const y of []) {}
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
const x = await 1;
export { x };
for await (const y of []) { }
//// [index.js]
const x = await 1;
export { x };
for await (const y of []) { }
