//// [tests/cases/compiler/moduleResolutionWithSuffixes_one.ts] ////

//// [index.ts]
import { ios } from "./foo";
//// [foo.ios.ts]
export function ios() {}
//// [foo.ts]
export function base() {}


//// [foo.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.base = base;
function base() { }
//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [foo.ios.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ios = ios;
function ios() { }
