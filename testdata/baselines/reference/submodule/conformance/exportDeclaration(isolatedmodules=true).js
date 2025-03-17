//// [tests/cases/conformance/externalModules/typeOnly/exportDeclaration.ts] ////

//// [a.ts]
class A {}
export type { A };

//// [b.ts]
import { A } from './a';
declare const a: A;
new A();

//// [c.ts]
import type { A } from './a';
export = A;

//// [d.ts]
import { A } from './a';
export = A;

//// [d.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [c.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
new A();
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
class A {
}
