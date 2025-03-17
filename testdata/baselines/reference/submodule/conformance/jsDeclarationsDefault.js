//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsDefault.ts] ////

//// [index1.js]
export default 12;

//// [index2.js]
export default function foo() {
    return foo;
}
export const x = foo;
export { foo as bar };

//// [index3.js]
export default class Foo {
    a = /** @type {Foo} */(null);
};
export const X = Foo;
export { Foo as Bar };

//// [index4.js]
import Fab from "./index3";
class Bar extends Fab {
    x = /** @type {Bar} */(null);
}
export default Bar;

//// [index5.js]
// merge type alias and const (OK)
export default 12;
/**
 * @typedef {string | number} default
 */

//// [index6.js]
// merge type alias and function (OK)
export default function func() {};
/**
 * @typedef {string | number} default
 */


//// [index1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = 12;
//// [index2.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.x = void 0;
exports.default = foo;
exports.bar = foo;
function foo() {
    return foo;
}
exports.x = foo;
//// [index3.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Bar = exports.X = void 0;
class Foo {
    a = (null);
}
exports.default = Foo;
exports.Bar = Foo;
;
exports.X = Foo;
//// [index4.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const index3_1 = require("./index3");
class Bar extends index3_1.default {
    x = (null);
}
exports.default = Bar;
//// [index5.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = 12;
//// [index6.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = func;
function func() { }
;
