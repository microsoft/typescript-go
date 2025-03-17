//// [tests/cases/compiler/resolutionCandidateFromPackageJsonField1.ts] ////

//// [package.json]
{
    "name": "@angular/core",
    "typings": "index.d.ts"
}

//// [index.ts]
export {};

//// [test.ts]
import "@angular/core";


//// [test.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
require("@angular/core");
//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
