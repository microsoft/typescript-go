//// [tests/cases/compiler/modulePreserve1.ts] ////

//// [a.ts]
export class A {}

//// [b.ts]
export = class B {}

//// [main.ts]
import { A } from "./a";
import B = require("./b");
export { A, B };


//// [main.js]
import { A } from "./a";
const B = require("./b");
export { A, B };
//// [b.js]
module.exports = class B {
};
//// [a.js]
export class A {
}
