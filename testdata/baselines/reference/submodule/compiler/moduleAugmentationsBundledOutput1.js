//// [tests/cases/compiler/moduleAugmentationsBundledOutput1.ts] ////

//// [m1.ts]
export class Cls {
}

//// [m2.ts]
import {Cls} from "./m1";
(<any>Cls.prototype).foo = function() { return 1; };
(<any>Cls.prototype).bar = function() { return "1"; };

declare module "./m1" {
    interface Cls {
        foo(): number;
    }
}

declare module "./m1" {
    interface Cls {
        bar(): string;
    }
}

//// [m3.ts]
export class C1 { x: number }
export class C2 { x: string }

//// [m4.ts]
import {Cls} from "./m1";
import {C1, C2} from "./m3";
(<any>Cls.prototype).baz1 = function() { return undefined };
(<any>Cls.prototype).baz2 = function() { return undefined };

declare module "./m1" {
    interface Cls {
        baz1(): C1;
    }
}

declare module "./m1" {
    interface Cls {
        baz2(): C2;
    }
}

//// [test.ts]
import { Cls } from "./m1";
import "m2";
import "m4";
let c: Cls;
c.foo().toExponential();
c.bar().toLowerCase();
c.baz1().x.toExponential();
c.baz2().x.toLowerCase();


//// [test.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
require("m2");
require("m4");
let c;
c.foo().toExponential();
c.bar().toLowerCase();
c.baz1().x.toExponential();
c.baz2().x.toLowerCase();
//// [m4.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const m1_1 = require("./m1");
m1_1.Cls.prototype.baz1 = function () { return undefined; };
m1_1.Cls.prototype.baz2 = function () { return undefined; };
//// [m3.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.C2 = exports.C1 = void 0;
class C1 {
    x;
}
exports.C1 = C1;
class C2 {
    x;
}
exports.C2 = C2;
//// [m2.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const m1_1 = require("./m1");
m1_1.Cls.prototype.foo = function () { return 1; };
m1_1.Cls.prototype.bar = function () { return "1"; };
//// [m1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Cls = void 0;
class Cls {
}
exports.Cls = Cls;
