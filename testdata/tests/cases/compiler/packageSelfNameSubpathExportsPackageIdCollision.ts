// @module: nodenext
// @declaration: true
// @declarationDir: /dist
// @rootDir: /
// @strict: true
// @traceResolution: true

// @filename: /dist/get.d.ts
export declare function get(object: any, path: string, defaultValue?: any): any;

// @filename: /dist/set.d.ts
export declare function set<T>(object: T, path: string, value: any): T;

// @filename: /get.ts
export function get(object: any, path: string, defaultValue?: any): any {
    return defaultValue;
}

// @filename: /set.ts
export function set<T>(object: T, path: string, value: any): T {
    return object;
}

// @filename: /main.ts
import { get } from "self/get";
import { set } from "self/set";

declare const obj: { a: { b: number } };
get(obj, "a.b");
set(obj, "a.b");

// @filename: /package.json
{
    "name": "self",
    "version": "1.0.0",
    "exports": {
        "./get": "./dist/get.d.ts",
        "./set": "./dist/set.d.ts"
    }
}
