//// [tests/cases/compiler/jsDeclarationEmitDoesNotRenameImport.ts] ////

//// [Test.js]
/** @module test/Test */
class Test {}
export default Test;
//// [Test.js]
/** @module Test */
class Test {}
export default Test;
//// [index.js]
import Test from './test/Test.js'

/**
 * @typedef {Object} Options
 * @property {typeof import("./Test.js").default} [test]
 */

class X extends Test {
    /**
     * @param {Options} options
     */
    constructor(options) {
        super();
        if (options.test) {
            this.test = new options.test();
        }
    }
}

export default X;


//// [Test.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
class Test {
}
exports.default = Test;
//// [Test.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
class Test {
}
exports.default = Test;
//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const Test_js_1 = require("./test/Test.js");
class X extends Test_js_1.default {
    constructor(options) {
        super();
        if (options.test) {
            this.test = new options.test();
        }
    }
}
exports.default = X;
