//// [tests/cases/compiler/outModuleConcatES6.ts] ////

//// [a.ts]
export class A { }

//// [b.ts]
import {A} from "./ref/a";
export class B extends A { }

//// [b.js]
import { A } from "./ref/a";
export class B extends A {
}
//// [a.js]
export class A {
}
