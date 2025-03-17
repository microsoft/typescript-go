//// [tests/cases/compiler/keepImportsInDts3.ts] ////

//// [test.ts]
export {};
//// [main.ts]
import "test"


//// [main.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
require("test");
//// [test.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
