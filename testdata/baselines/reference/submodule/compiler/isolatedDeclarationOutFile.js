//// [tests/cases/compiler/isolatedDeclarationOutFile.ts] ////

//// [a.ts]
export class A {
    toUpper(msg: string): string {
        return msg.toUpperCase();
    }
}

//// [b.ts]
import { A } from "./a";

export class B extends A {
    toFixed(n: number): string {
        return n.toFixed(6);
    }
}

export function makeB(): A {
    return new B();
}


//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.B = void 0;
exports.makeB = makeB;
const a_1 = require("./a");
class B extends a_1.A {
    toFixed(n) {
        return n.toFixed(6);
    }
}
exports.B = B;
function makeB() {
    return new B();
}
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.A = void 0;
class A {
    toUpper(msg) {
        return msg.toUpperCase();
    }
}
exports.A = A;
