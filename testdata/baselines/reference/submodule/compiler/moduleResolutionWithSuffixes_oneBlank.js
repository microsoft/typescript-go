//// [tests/cases/compiler/moduleResolutionWithSuffixes_oneBlank.ts] ////

//// [index.ts]
import { base } from "./foo";
//// [foo.ts]
export function base() {}


//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [foo.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.base = base;
function base() { }
