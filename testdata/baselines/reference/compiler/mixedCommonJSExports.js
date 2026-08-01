//// [tests/cases/compiler/mixedCommonJSExports.ts] ////

//// [requires.d.ts]
declare var module: { exports: any };
declare function require(name: string): any;

//// [mod.js]
module.exports = function main() {};
/** @param {number} value */
module.exports.helper = function helper(value) {};

//// [plugin.js]
class Plugin {}
module.exports = Plugin;
module.exports.helper = function helper() {};

//// [defined.js]
module.exports = class Defined {};
Object.defineProperty(module.exports, "value", { value: 1 });

//// [arrow.js]
const arrow = () => {};
module.exports = arrow;
module.exports.helper = function helper() {};

//// [index.js]
const mod = require("./mod");
const Plugin = require("./plugin");
const Defined = require("./defined");
const arrow = require("./arrow");

mod();
mod.helper(1);
Plugin.helper();
Defined.value.toFixed();
arrow.helper();




//// [mod.d.ts]
export = main;
declare function main(): void;
declare namespace main {
    export var helper: (value: number) => void;
}
//// [plugin.d.ts]
export = Plugin;
declare namespace Plugin {
    export var helper: () => void;
}
declare class Plugin {
}
//// [defined.d.ts]
export = Defined;
declare class Defined {
}
declare namespace Defined {
    export var value: number;
}
//// [arrow.d.ts]
export = _exports;
declare function _exports(): void;
declare namespace _exports {
    export var helper: () => void;
}
declare const arrow: () => void;
//// [index.d.ts]
export {};
