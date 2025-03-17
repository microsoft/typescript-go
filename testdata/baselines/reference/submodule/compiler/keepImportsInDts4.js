//// [tests/cases/compiler/keepImportsInDts4.ts] ////

//// [test.ts]
export {};
//// [main.ts]
import "./folder/test"


//// [test.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [main.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
require("./folder/test");
