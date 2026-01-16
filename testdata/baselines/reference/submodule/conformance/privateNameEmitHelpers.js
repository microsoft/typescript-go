//// [tests/cases/conformance/classes/members/privateNames/privateNameEmitHelpers.ts] ////

//// [main.ts]
export class C {
    #a = 1;
    #b() { this.#c = 42; }
    set #c(v: number) { this.#a += v; }
}

//// [index.d.ts]
// these are pre-TS4.3 versions of emit helpers, which only supported private instance fields
export declare function __classPrivateFieldGet<T extends object, V>(receiver: T, state: any): V;
export declare function __classPrivateFieldSet<T extends object, V>(receiver: T, state: any, value: V): V;


//// [main.js]
import { __classPrivateFieldGet, __classPrivateFieldSet } from "tslib";
var _C_a;
export class C {
    constructor() {
        _C_a.set(this, 1);
    }
    #b() { this.#c = 42; }
    set #c(v) { __classPrivateFieldSet(this, _C_a, __classPrivateFieldGet(this, _C_a, "f") + v, "f"); }
}
_C_a = new WeakMap();
