// @noTypesAndSymbols: true
// @target: es2015
// @module: commonjs
// @importHelpers: true

// @filename: /node_modules/tslib/package.json
{
    "name": "tslib",
    "main": "tslib.js",
    "typings": "tslib.d.ts"
}

// @filename: /node_modules/tslib/tslib.d.ts
export declare function __awaiter(thisArg: any, _arguments: any, P: Function, generator: Function): any;

// @filename: /main.ts
export async function doStuff() {
    return 1;
}
