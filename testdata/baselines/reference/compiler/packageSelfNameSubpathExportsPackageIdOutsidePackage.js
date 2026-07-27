//// [tests/cases/compiler/packageSelfNameSubpathExportsPackageIdOutsidePackage.ts] ////

//// [get.d.ts]
export declare function get(object: any, path: string, defaultValue?: any): any;

//// [set.d.ts]
export declare function set<T>(object: T, path: string, value: any): T;

//// [get.ts]
export function get(object: any, path: string, defaultValue?: any): any {
    return defaultValue;
}

//// [set.ts]
export function set<T>(object: T, path: string, value: any): T {
    return object;
}

//// [main.ts]
import { get } from "self/get";
import { set } from "self/set";

declare const obj: { a: { b: number } };
get(obj, "a.b");
set(obj, "a.b");

//// [package.json]
{
    "name": "self",
    "version": "1.0.0",
    "exports": {
        "./get": "./dist/get.d.ts",
        "./set": "./dist/set.d.ts"
    }
}


//// [get.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.get = get;
function get(object, path, defaultValue) {
    return defaultValue;
}
//// [main.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const get_1 = require("self/get");
const set_1 = require("self/set");
(0, get_1.get)(obj, "a.b");
(0, set_1.set)(obj, "a.b");


//// [main.d.ts]
export {};
