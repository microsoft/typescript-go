//// [tests/cases/compiler/emitBundleWithPrologueDirectives1.ts] ////

//// [test.ts]
/* Detached Comment */

// Class Doo Comment
export class Doo {}
class Scooby extends Doo {}

//// [test.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Doo = void 0;
class Doo {
}
exports.Doo = Doo;
class Scooby extends Doo {
}
