//// [tests/cases/conformance/externalModules/typeOnly/implementsClause.ts] ////

//// [types.ts]
export interface Component {}

//// [ns.ts]
import type * as types from './types';
export { types };

//// [index.ts]
import type * as types from './types';
import * as nestedNamespace from './ns';

class C implements types.Component {}
class D implements nestedNamespace.types.Component {}


//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
class C {
}
class D {
}
//// [ns.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [types.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
