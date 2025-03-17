//// [tests/cases/conformance/externalModules/typeOnly/chained.ts] ////

//// [a.ts]
class A { a!: string }
export type { A as B };
export type Z = A;

//// [b.ts]
import { Z as Y } from './a';
export { B as C } from './a';

//// [c.ts]
import type { C } from './b';
export { C as D };

//// [d.ts]
import { D } from './c';
new D();
const d: D = {};


//// [d.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
new D();
const d = {};
//// [c.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
class A {
    a;
}
