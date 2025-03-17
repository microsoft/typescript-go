//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsClassStatic2.ts] ////

//// [Foo.js]
class Base {
  static foo = "";
}
export class Foo extends Base {}
Foo.foo = "foo";

//// [Bar.ts]
import { Foo } from "./Foo.js";

class Bar extends Foo {}
Bar.foo = "foo";


//// [Bar.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const Foo_js_1 = require("./Foo.js");
class Bar extends Foo_js_1.Foo {
}
Bar.foo = "foo";
//// [Foo.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Foo = void 0;
class Base {
    static foo = "";
}
class Foo extends Base {
}
exports.Foo = Foo;
Foo.foo = "foo";
