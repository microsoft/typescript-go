//// [tests/cases/conformance/es6/moduleExportsAmd/outFilerootDirModuleNamesAmd.ts] ////

//// [a.ts]
import foo from "./b";
export default class Foo {}
foo();

//// [b.ts]
import Foo from "./a";
export default function foo() { new Foo(); }

// https://github.com/microsoft/TypeScript/issues/37429
import("./a");

//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const b_1 = require("./b");
class Foo {
}
exports.default = Foo;
(0, b_1.default)();
//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = foo;
const a_1 = require("./a");
function foo() { new a_1.default(); }
Promise.resolve().then(() => require("./a"));
