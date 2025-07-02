//// [tests/cases/compiler/declarationEmitReexportedSymlinkReference2.ts] ////

//// [index.d.ts]
export * from './types';
//// [types.d.ts]
export declare type A = {
    id: string;
};
export declare type B = {
    id: number;
};
export declare type IdType = A | B;
export declare class MetadataAccessor<T, D extends IdType = IdType> {
    readonly key: string;
    private constructor();
    toString(): string;
    static create<T, D extends IdType = IdType>(key: string): MetadataAccessor<T, D>;
}
//// [package.json]
{
    "name": "@raymondfeng/pkg1",
    "version": "1.0.0",
    "description": "",
    "main": "dist/index.js",
    "typings": "dist/index.d.ts"
}
//// [index.d.ts]
import "./secondary";
export * from './types';
//// [types.d.ts]
export {MetadataAccessor} from '@raymondfeng/pkg1';
//// [secondary.d.ts]
export {IdType} from '@raymondfeng/pkg1';
//// [package.json]
{
    "name": "@raymondfeng/pkg2",
    "version": "1.0.0",
    "description": "",
    "main": "dist/index.js",
    "typings": "dist/index.d.ts"
}
//// [index.ts]
export * from './keys';
//// [keys.ts]
import {MetadataAccessor} from "@raymondfeng/pkg2";

export const ADMIN = MetadataAccessor.create<boolean>('1');

//// [keys.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ADMIN = void 0;
const pkg2_1 = require("@raymondfeng/pkg2");
exports.ADMIN = pkg2_1.MetadataAccessor.create('1');
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
var __exportStar = (this && this.__exportStar) || function(m, exports) {
    for (var p in m) if (p !== "default" && !Object.prototype.hasOwnProperty.call(exports, p)) __createBinding(exports, m, p);
};
Object.defineProperty(exports, "__esModule", { value: true });
__exportStar(require("./keys"), exports);


//// [keys.d.ts]
import { MetadataAccessor } from "@raymondfeng/pkg2";
export declare const ADMIN: MetadataAccessor<boolean, import("@raymondfeng/pkg2").IdType>;
//// [index.d.ts]
export * from './keys';


//// [DtsFileErrors]


monorepo/pkg3/dist/keys.d.ts(2,83): error TS2694: Namespace '"monorepo/pkg2/dist/index"' has no exported member 'IdType'.


==== monorepo/pkg3/tsconfig.json (0 errors) ====
    {
        "compilerOptions": {
          "outDir": "dist",
          "rootDir": "src",
          "target": "es5",
          "module": "commonjs",
          "strict": true,
          "esModuleInterop": true,
          "declaration": true
        }
    }
    
==== monorepo/pkg3/dist/index.d.ts (0 errors) ====
    export * from './keys';
    
==== monorepo/pkg3/dist/keys.d.ts (1 errors) ====
    import { MetadataAccessor } from "@raymondfeng/pkg2";
    export declare const ADMIN: MetadataAccessor<boolean, import("@raymondfeng/pkg2").IdType>;
                                                                                      ~~~~~~
!!! error TS2694: Namespace '"monorepo/pkg2/dist/index"' has no exported member 'IdType'.
    
==== monorepo/pkg1/dist/index.d.ts (0 errors) ====
    export * from './types';
==== monorepo/pkg1/dist/types.d.ts (0 errors) ====
    export declare type A = {
        id: string;
    };
    export declare type B = {
        id: number;
    };
    export declare type IdType = A | B;
    export declare class MetadataAccessor<T, D extends IdType = IdType> {
        readonly key: string;
        private constructor();
        toString(): string;
        static create<T, D extends IdType = IdType>(key: string): MetadataAccessor<T, D>;
    }
==== monorepo/pkg1/package.json (0 errors) ====
    {
        "name": "@raymondfeng/pkg1",
        "version": "1.0.0",
        "description": "",
        "main": "dist/index.js",
        "typings": "dist/index.d.ts"
    }
==== monorepo/pkg2/dist/index.d.ts (0 errors) ====
    import "./secondary";
    export * from './types';
==== monorepo/pkg2/dist/types.d.ts (0 errors) ====
    export {MetadataAccessor} from '@raymondfeng/pkg1';
==== monorepo/pkg2/dist/secondary.d.ts (0 errors) ====
    export {IdType} from '@raymondfeng/pkg1';
==== monorepo/pkg2/package.json (0 errors) ====
    {
        "name": "@raymondfeng/pkg2",
        "version": "1.0.0",
        "description": "",
        "main": "dist/index.js",
        "typings": "dist/index.d.ts"
    }