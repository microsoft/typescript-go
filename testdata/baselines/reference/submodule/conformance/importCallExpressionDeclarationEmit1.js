//// [tests/cases/conformance/dynamicImport/importCallExpressionDeclarationEmit1.ts] ////

//// [importCallExpressionDeclarationEmit1.ts]
declare function getSpecifier(): string;
declare var whatToLoad: boolean;
declare const directory: string;
declare const moduleFile: number;

import(getSpecifier());

var p0 = import(`${directory}\\${moduleFile}`);
var p1 = import(getSpecifier());
const p2 = import(whatToLoad ? getSpecifier() : "defaulPath")

function returnDynamicLoad(path: string) {
    return import(path);
}

//// [importCallExpressionDeclarationEmit1.js]
"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Promise.resolve(`${getSpecifier()}`).then(s => __importStar(require(s)));
var p0 = Promise.resolve(`${`${directory}\\${moduleFile}`}`).then(s => __importStar(require(s)));
var p1 = Promise.resolve(`${getSpecifier()}`).then(s => __importStar(require(s)));
const p2 = Promise.resolve(`${whatToLoad ? getSpecifier() : "defaulPath"}`).then(s => __importStar(require(s)));
function returnDynamicLoad(path) {
    return Promise.resolve(`${path}`).then(s => __importStar(require(s)));
}


//// [importCallExpressionDeclarationEmit1.d.ts]
function getSpecifier(): string;
var whatToLoad: boolean;
const directory: string;
const moduleFile: number;
var p0: Promise<any>;
var p1: Promise<any>;
const p2: Promise<any>;
function returnDynamicLoad(path: string): Promise<any>;


//// [DtsFileErrors]


importCallExpressionDeclarationEmit1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== importCallExpressionDeclarationEmit1.d.ts (1 errors) ====
    function getSpecifier(): string;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    var whatToLoad: boolean;
    const directory: string;
    const moduleFile: number;
    var p0: Promise<any>;
    var p1: Promise<any>;
    const p2: Promise<any>;
    function returnDynamicLoad(path: string): Promise<any>;
    