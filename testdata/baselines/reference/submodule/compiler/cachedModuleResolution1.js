//// [tests/cases/compiler/cachedModuleResolution1.ts] ////

//// [foo.d.ts]
export declare let x: number

//// [app.ts]
import {x} from "foo";

//// [lib.ts]
import {x} from "foo";


//// [lib.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [app.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
