//// [tests/cases/conformance/externalModules/typeOnly/grammarErrors.ts] ////

//// [a.ts]
export default class A {}
export class B {}
export class C {}

//// [b.ts]
import type A, { B, C } from './a';

//// [a.js]
import type A from './a';
export type { A };

//// [c.ts]
namespace ns {
  export class Foo {}
}
import type Foo = ns.Foo;


//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.C = exports.B = void 0;
class A {
}
exports.default = A;
class B {
}
exports.B = B;
class C {
}
exports.C = C;
//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [c.js]
var ns;
(function (ns) {
    class Foo {
    }
    ns.Foo = Foo;
})(ns || (ns = {}));
