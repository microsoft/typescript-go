//// [tests/cases/conformance/classes/members/privateNames/privateNameFieldParenthesisLeftAssignment.ts] ////

//// [privateNameFieldParenthesisLeftAssignment.ts]
class Foo {
    #p: number;

    constructor(value: number) {
        this.#p = value;
    }

    t1(p: number) {
        (this.#p as number) = p;
    }

    t2(p: number) {
        (((this.#p as number))) = p;
    }

    t3(p: number) {
        (this.#p) = p;
    }

    t4(p: number) {
        (((this.#p))) = p;
    }
}


//// [privateNameFieldParenthesisLeftAssignment.js]
var __classPrivateFieldSet = (this && this.__classPrivateFieldSet) || function (receiver, state, value, kind, f) {
    if (kind === "m") throw new TypeError("Private method is not writable");
    if (kind === "a" && !f) throw new TypeError("Private accessor was defined without a setter");
    if (typeof state === "function" ? receiver !== state || !f : !state.has(receiver)) throw new TypeError("Cannot write private member to an object whose class did not declare it");
    return (kind === "a" ? f.call(receiver, value) : f ? f.value = value : state.set(receiver, value)), value;
};
var __classPrivateFieldGet = (this && this.__classPrivateFieldGet) || function (receiver, state, kind, f) {
    if (kind === "a" && !f) throw new TypeError("Private accessor was defined without a getter");
    if (typeof state === "function" ? receiver !== state || !f : !state.has(receiver)) throw new TypeError("Cannot read private member from an object whose class did not declare it");
    return kind === "m" ? f : kind === "a" ? f.call(receiver) : f ? f.value : state.get(receiver);
};
var _Foo_p;
class Foo {
    constructor(value) {
        _Foo_p.set(this, void 0);
        __classPrivateFieldSet(this, _Foo_p, value, "f");
    }
    t1(p) {
        __classPrivateFieldGet(this, _Foo_p, "f") = p;
    }
    t2(p) {
        __classPrivateFieldGet(this, _Foo_p, "f") = p;
    }
    t3(p) {
        __classPrivateFieldSet(this, _Foo_p, p, "f");
    }
    t4(p) {
        __classPrivateFieldSet(this, _Foo_p, p, "f");
    }
}
_Foo_p = new WeakMap();
