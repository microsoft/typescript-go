//// [tests/cases/compiler/checkExternalEmitHelpersCrash.ts] ////

//// [package.json]
{
    "name": "tslib",
    "main": "tslib.js",
    "typings": "tslib.d.ts"
}

//// [tslib.d.ts]
export declare function __awaiter(thisArg: any, _arguments: any, P: Function, generator: Function): any;

//// [main.ts]
export async function doStuff() {
    return 1;
}


//// [main.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.doStuff = doStuff;
const tslib_1 = require("tslib");
function doStuff() {
    return tslib_1.__awaiter(this, void 0, void 0, function* () {
        return 1;
    });
}
