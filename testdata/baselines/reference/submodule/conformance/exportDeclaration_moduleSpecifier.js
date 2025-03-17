//// [tests/cases/conformance/externalModules/typeOnly/exportDeclaration_moduleSpecifier.ts] ////

//// [a.ts]
export class A {}

//// [b.ts]
export type { A } from './a';

//// [c.ts]
import { A } from './b';
declare const a: A;
new A();


//// [c.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
new A();
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
