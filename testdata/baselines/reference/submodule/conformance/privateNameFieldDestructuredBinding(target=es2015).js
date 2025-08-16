//// [tests/cases/conformance/classes/members/privateNames/privateNameFieldDestructuredBinding.ts] ////

//// [privateNameFieldDestructuredBinding.ts]
class A {
    #field = 1;
    otherObject = new A();
    testObject() {
        return { x: 10, y: 6 };
    }
    testArray() {
        return [10, 11];
    }
    constructor() {
        let y: number;
        ({ x: this.#field, y } = this.testObject());
        ([this.#field, y] = this.testArray());
        ({ a: this.#field, b: [this.#field] } = { a: 1, b: [2] });
        [this.#field, [this.#field]] = [1, [2]];
        ({ a: this.#field = 1, b: [this.#field = 1] } = { b: [] });
        [this.#field = 2] = [];
        [this.otherObject.#field = 2] = [];
    }
    static test(_a: A) {
        [_a.#field] = [2];
    }
}


//// [privateNameFieldDestructuredBinding.js]
var __classPrivateFieldGet = (this && this.__classPrivateFieldGet) || function (receiver, state, kind, f) {
    if (kind === "a" && !f) throw new TypeError("Private accessor was defined without a getter");
    if (typeof state === "function" ? receiver !== state || !f : !state.has(receiver)) throw new TypeError("Cannot read private member from an object whose class did not declare it");
    return kind === "m" ? f : kind === "a" ? f.call(receiver) : f ? f.value : state.get(receiver);
};
var __classPrivateFieldSet = (this && this.__classPrivateFieldSet) || function (receiver, state, value, kind, f) {
    if (kind === "m") throw new TypeError("Private method is not writable");
    if (kind === "a" && !f) throw new TypeError("Private accessor was defined without a setter");
    if (typeof state === "function" ? receiver !== state || !f : !state.has(receiver)) throw new TypeError("Cannot write private member to an object whose class did not declare it");
    return (kind === "a" ? f.call(receiver, value) : f ? f.value = value : state.set(receiver, value)), value;
};
var _A_field;
class A {
    otherObject = new A();
    testObject() {
        return { x: 10, y: 6 };
    }
    testArray() {
        return [10, 11];
    }
    constructor() {
        _A_field.set(this, 1);
        let y;
        ({ x: __classPrivateFieldGet(this, _A_field, "f"), y } = this.testObject());
        ([__classPrivateFieldGet(this, _A_field, "f"), y] = this.testArray());
        ({ a: __classPrivateFieldGet(this, _A_field, "f"), b: [__classPrivateFieldGet(this, _A_field, "f")] } = { a: 1, b: [2] });
        [__classPrivateFieldGet(this, _A_field, "f"), [__classPrivateFieldGet(this, _A_field, "f")]] = [1, [2]];
        ({ a: __classPrivateFieldSet(this, _A_field, 1, "f"), b: [__classPrivateFieldSet(this, _A_field, 1, "f")] } = { b: [] });
        [__classPrivateFieldSet(this, _A_field, 2, "f")] = [];
        [__classPrivateFieldSet(this.otherObject, _A_field, 2, "f")] = [];
    }
    static test(_a) {
        [__classPrivateFieldGet(_a, _A_field, "f")] = [2];
    }
}
_A_field = new WeakMap();
