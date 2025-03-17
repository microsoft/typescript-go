//// [tests/cases/compiler/importNonExportedMember2.ts] ////

//// [a.ts]
export {}
interface Foo {}

//// [b.ts]
import { Foo } from './a';


//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
