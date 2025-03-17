//// [tests/cases/conformance/externalModules/typeOnly/renamed.ts] ////

//// [a.ts]
class A { a!: string }
export type { A as B };

//// [b.ts]
export type { B as C } from './a';

//// [c.ts]
import type { C as D } from './b';
const d: D = {};


//// [c.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const d = {};
//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
class A {
    a;
}
