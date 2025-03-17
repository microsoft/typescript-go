//// [tests/cases/conformance/moduleResolution/bundler/bundlerImportESM.ts] ////

//// [esm.mts]
export const esm = 0;

//// [not-actually-cjs.cts]
import { esm } from "./esm.mjs";

//// [package.json]
{ "type": "commonjs" }

//// [still-not-cjs.ts]
import { esm } from "./esm.mjs";


//// [still-not-cjs.js]
export {};
//// [not-actually-cjs.cjs]
export {};
//// [esm.mjs]
export const esm = 0;
