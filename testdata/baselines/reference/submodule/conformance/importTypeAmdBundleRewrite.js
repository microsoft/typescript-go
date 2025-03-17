//// [tests/cases/conformance/types/import/importTypeAmdBundleRewrite.ts] ////

//// [c.ts]
export interface Foo {
    x: 12;
}
//// [inner.ts]
const c: import("./b/c").Foo = {x: 12};
export {c};

//// [index.ts]
const d: typeof import("./a/inner")["c"] = {x: 12};
export {d};


//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.d = void 0;
const d = { x: 12 };
exports.d = d;
//// [inner.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.c = void 0;
const c = { x: 12 };
exports.c = c;
//// [c.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
