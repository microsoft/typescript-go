// @module: nodenext
// @declaration: true
// @declarationDir: /very/long/package/dist
// @rootDir: /src
// @strict: true
// @traceResolution: true

// @filename: /very/long/package/dist/get.d.ts
export declare function get(object: any, path: string, defaultValue?: any): any;

// @filename: /very/long/package/dist/set.d.ts
export declare function set<T>(object: T, path: string, value: any): T;

// @filename: /src/get.ts
export function get(object: any, path: string, defaultValue?: any): any {
    return defaultValue;
}

// @filename: /src/set.ts
export function set<T>(object: T, path: string, value: any): T {
    return object;
}

// @filename: /very/long/package/main.ts
import { get } from "self/get";
import { set } from "self/set";

declare const obj: { a: { b: number } };
get(obj, "a.b");
set(obj, "a.b");

// @filename: /very/long/package/package.json
{
    "name": "self",
    "version": "1.0.0",
    "exports": {
        "./get": "./dist/get.d.ts",
        "./set": "./dist/set.d.ts"
    }
}
