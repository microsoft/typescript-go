//// [tests/cases/compiler/keepImportsInDts2.ts] ////

//// [test.ts]
export {};
//// [main.ts]
import "./folder/test"

//// [main.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
require("./folder/test");
//// [test.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
