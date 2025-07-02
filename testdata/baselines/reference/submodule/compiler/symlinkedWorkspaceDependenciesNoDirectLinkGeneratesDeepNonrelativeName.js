//// [tests/cases/compiler/symlinkedWorkspaceDependenciesNoDirectLinkGeneratesDeepNonrelativeName.ts] ////

//// [foo.d.ts]
export declare class Foo {
    private f: any;
}
//// [index.d.ts]
import { Foo } from "./foo.js";
export function create(): Foo;
//// [package.json]
{
    "name": "package-a",
    "version": "0.0.1",
    "exports": {
        ".": "./index.js",
        "./cls": "./foo.js"
    }
}
//// [package.json]
{
    "private": true,
    "dependencies": {
        "package-a": "file:../packageA"
    }
}
//// [index.d.ts]
import { create } from "package-a";
export declare function invoke(): ReturnType<typeof create>;
//// [package.json]
{
    "private": true,
    "dependencies": {
        "package-b": "file:../packageB",
        "package-a": "file:../packageA"
    }
}
//// [index.ts]
import * as pkg from "package-b";

export const a = pkg.invoke();

//// [index.js]
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
Object.defineProperty(exports, "__esModule", { value: true });
exports.a = void 0;
const pkg = __importStar(require("package-b"));
exports.a = pkg.invoke();


//// [index.d.ts]
export declare const a: import("package-a").Foo;


//// [DtsFileErrors]


workspace/packageC/index.d.ts(1,45): error TS2694: Namespace '"workspace/packageA/index"' has no exported member 'Foo'.


==== workspace/packageA/foo.d.ts (0 errors) ====
    export declare class Foo {
        private f: any;
    }
==== workspace/packageA/index.d.ts (0 errors) ====
    import { Foo } from "./foo.js";
    export function create(): Foo;
==== workspace/packageA/package.json (0 errors) ====
    {
        "name": "package-a",
        "version": "0.0.1",
        "exports": {
            ".": "./index.js",
            "./cls": "./foo.js"
        }
    }
==== workspace/packageB/package.json (0 errors) ====
    {
        "private": true,
        "dependencies": {
            "package-a": "file:../packageA"
        }
    }
==== workspace/packageB/index.d.ts (0 errors) ====
    import { create } from "package-a";
    export declare function invoke(): ReturnType<typeof create>;
==== workspace/packageC/package.json (0 errors) ====
    {
        "private": true,
        "dependencies": {
            "package-b": "file:../packageB",
            "package-a": "file:../packageA"
        }
    }
==== workspace/packageC/index.d.ts (1 errors) ====
    export declare const a: import("package-a").Foo;
                                                ~~~
!!! error TS2694: Namespace '"workspace/packageA/index"' has no exported member 'Foo'.
    