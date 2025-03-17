//// [tests/cases/compiler/importNonExportedMember1.ts] ////

//// [a.ts]
declare function foo(): any
declare function bar(): any;
export { foo };

//// [b.ts]
import { bar } from "./a";


//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.foo = foo;
