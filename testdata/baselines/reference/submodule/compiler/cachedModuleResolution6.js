//// [tests/cases/compiler/cachedModuleResolution6.ts] ////

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
