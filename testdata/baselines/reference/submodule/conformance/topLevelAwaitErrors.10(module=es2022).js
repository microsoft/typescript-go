//// [tests/cases/conformance/externalModules/topLevelAwaitErrors.10.ts] ////

//// [index.ts]
// await disallowed in alias of named import
import { await as await } from "./other";

//// [other.ts]
declare const _await: any;
export { _await as await };


//// [index.js]
export {};
//// [other.js]
export { _await as await };
