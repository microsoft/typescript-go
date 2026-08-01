//// [tests/cases/compiler/commonJSExportAssignmentVariableInitializers.ts] ////

//// [requires.d.ts]
declare var module: { exports: any };
declare function require(name: string): any;

//// [uninitialized.js]
var uninitialized;
module.exports = uninitialized;

//// [object.js]
const object = { value: 1 };
module.exports = object;
module.exports.helper = function helper() {};

//// [index.js]
require("./uninitialized");
const object = require("./object");
object.helper();




//// [uninitialized.d.ts]
export = uninitialized;
declare var uninitialized: any;
//// [object.d.ts]
export = object;
export declare var helper: () => void;
declare const object: {
    value: number;
};
//// [index.d.ts]
export {};
