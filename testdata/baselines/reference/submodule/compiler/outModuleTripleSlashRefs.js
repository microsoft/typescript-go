//// [tests/cases/compiler/outModuleTripleSlashRefs.ts] ////

//// [a.ts]
/// <reference path="./b.ts" />
export class A {
	member: typeof GlobalFoo;
}

//// [b.ts]
/// <reference path="./c.d.ts" />
class Foo {
	member: Bar;
}
declare var GlobalFoo: Foo;

//// [c.d.ts]
/// <reference path="./d.d.ts" />
declare class Bar {
	member: Baz;
}

//// [d.d.ts]
declare class Baz {
	member: number;
}

//// [b.ts]
import {A} from "./ref/a";
export class B extends A { }


//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.B = void 0;
const a_1 = require("./ref/a");
class B extends a_1.A {
}
exports.B = B;
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.A = void 0;
class A {
    member;
}
exports.A = A;
//// [b.js]
class Foo {
    member;
}
