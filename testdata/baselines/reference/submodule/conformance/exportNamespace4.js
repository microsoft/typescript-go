//// [tests/cases/conformance/externalModules/typeOnly/exportNamespace4.ts] ////

//// [a.ts]
export class A {}

//// [b.ts]
export type * from './a';

//// [c.ts]
export type * as ns from './a';

//// [d.ts]
import { A } from './b';
A;

//// [e.ts]
import { ns } from './c';
ns.A;


//// [e.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
ns.A;
//// [d.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
A;
//// [c.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.A = void 0;
class A {
}
exports.A = A;
