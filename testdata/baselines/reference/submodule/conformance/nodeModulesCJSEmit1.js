//// [tests/cases/conformance/node/nodeModulesCJSEmit1.ts] ////

//// [1.cjs]
module.exports = {};

//// [2.cjs]
exports.foo = 0;

//// [3.cjs]
import "foo";
exports.foo = {};

//// [4.cjs]
;

//// [5.cjs]
import two from "./2.cjs";   // ok
import three from "./3.cjs"; // error
two.foo;
three.foo;


//// [1.cjs]
module.exports = {};
//// [2.cjs]
exports.foo = 0;
//// [3.cjs]
import "foo";
exports.foo = {};
//// [4.cjs]
;
//// [5.cjs]
import two from "./2.cjs";
import three from "./3.cjs";
two.foo;
three.foo;
