//// [tests/cases/compiler/nodeModuleReexportFromDottedPath.ts] ////

//// [index.d.ts]
export interface PrismaClientOptions {
  rejectOnNotFound?: any;
}

export class PrismaClient<T extends PrismaClientOptions = PrismaClientOptions> {
  private fetcher;
}

//// [index.d.ts]
export * from ".prisma/client";

//// [index.ts]
import { PrismaClient } from "@prisma/client";
declare const enhancePrisma: <TPrismaClientCtor>(client: TPrismaClientCtor) => TPrismaClientCtor & { enhanced: unknown };
const EnhancedPrisma = enhancePrisma(PrismaClient);
export default new EnhancedPrisma();


//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const client_1 = require("@prisma/client");
const EnhancedPrisma = enhancePrisma(client_1.PrismaClient);
exports.default = new EnhancedPrisma();


//// [index.d.ts]
import { PrismaClient } from "@prisma/client";
const _default: PrismaClient<import(".prisma/client").PrismaClientOptions>;
export default _default;


//// [DtsFileErrors]


/index.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== /node_modules/.prisma/client/index.d.ts (0 errors) ====
    export interface PrismaClientOptions {
      rejectOnNotFound?: any;
    }
    
    export class PrismaClient<T extends PrismaClientOptions = PrismaClientOptions> {
      private fetcher;
    }
    
==== /node_modules/@prisma/client/index.d.ts (0 errors) ====
    export * from ".prisma/client";
    
==== /index.d.ts (1 errors) ====
    import { PrismaClient } from "@prisma/client";
    const _default: PrismaClient<import(".prisma/client").PrismaClientOptions>;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    