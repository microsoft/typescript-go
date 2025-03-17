//// [tests/cases/compiler/pathMappingBasedModuleResolution_rootImport_aliasWithRoot.ts] ////

//// [foo.ts]
export function foo() {}

//// [bar.js]
export function bar() {}

//// [a.ts]
import { foo } from "/foo";
import { bar } from "/bar";


//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [bar.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.bar = bar;
function bar() { }
//// [foo.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.foo = foo;
function foo() { }
