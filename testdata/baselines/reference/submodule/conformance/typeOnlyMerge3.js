//// [tests/cases/conformance/externalModules/typeOnlyMerge3.ts] ////

//// [a.ts]
function A() {}
export type { A };

//// [b.ts]
import { A } from "./a";
namespace A {
  export const displayName = "A";
}
export { A };

//// [c.ts]
import { A } from "./b";
A;
A.displayName;
A();


//// [c.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const b_1 = require("./b");
A;
A.displayName;
A();
//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.A = void 0;
var A;
(function (A) {
    A.displayName = "A";
})(A || (exports.A = A = {}));
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
function A() { }


!!!! File c.js differs from original emit in noCheck emit
//// [c.js]
--- Expected	The full check baseline
+++ Actual	with noCheck set
@@ -1,6 +1,6 @@
 "use strict";
 Object.defineProperty(exports, "__esModule", { value: true });
 const b_1 = require("./b");
-A;
-A.displayName;
-A();
+b_1.A;
+b_1.A.displayName;
+(0, b_1.A)();
