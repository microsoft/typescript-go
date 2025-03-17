//// [tests/cases/conformance/es6/moduleExportsCommonjs/defaultExportsGetExportedCommonjs.ts] ////

//// [a.ts]
export default class Foo {}

//// [b.ts]
export default function foo() {}


//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = foo;
function foo() { }
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
class Foo {
}
exports.default = Foo;
