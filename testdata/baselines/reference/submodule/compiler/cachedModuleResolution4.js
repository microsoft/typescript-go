//// [tests/cases/compiler/cachedModuleResolution4.ts] ////

//// [foo.d.ts]
export declare let x: number

//// [lib.ts]
import {x} from "foo";

//// [app.ts]
import {x} from "foo";


//// [app.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [lib.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
