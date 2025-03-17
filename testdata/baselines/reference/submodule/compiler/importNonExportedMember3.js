//// [tests/cases/compiler/importNonExportedMember3.ts] ////

//// [a.ts]
export {}
interface Foo {}
interface Foo {}
namespace Foo {}

//// [b.ts]
import { Foo } from './a';


//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
