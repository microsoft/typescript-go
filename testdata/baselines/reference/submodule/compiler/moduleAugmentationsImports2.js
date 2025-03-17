//// [tests/cases/compiler/moduleAugmentationsImports2.ts] ////

//// [a.ts]
export class A {}

//// [b.ts]
export class B {x: number;}

//// [c.d.ts]
declare module "C" {
    class Cls {y: string; }
}

//// [d.ts]
/// <reference path="c.d.ts"/>

import {A} from "./a";
import {B} from "./b";

A.prototype.getB = function () { return undefined; }

declare module "./a" {
    interface A {
        getB(): B;
    }
}

//// [e.ts]
import {A} from "./a";
import {Cls} from "C";

A.prototype.getCls = function () { return undefined; }

declare module "./a" {
    interface A {
        getCls(): Cls;
    }
}

//// [main.ts]
import {A} from "./a";
import "d";
import "e";

let a: A;
let b = a.getB().x.toFixed();
let c = a.getCls().y.toLowerCase();


//// [main.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
require("d");
require("e");
let a;
let b = a.getB().x.toFixed();
let c = a.getCls().y.toLowerCase();
//// [e.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const a_1 = require("./a");
a_1.A.prototype.getCls = function () { return undefined; };
//// [d.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const a_1 = require("./a");
a_1.A.prototype.getB = function () { return undefined; };
//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.B = void 0;
class B {
    x;
}
exports.B = B;
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.A = void 0;
class A {
}
exports.A = A;
