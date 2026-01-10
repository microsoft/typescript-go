//// [tests/cases/conformance/classes/members/privateNames/privateNameFieldUnaryMutation.ts] ////

//// [privateNameFieldUnaryMutation.ts]
class C {
    #test: number = 24;
    constructor() {
        this.#test++;
        this.#test--;
        ++this.#test;
        --this.#test;
        const a = this.#test++;
        const b = this.#test--;
        const c = ++this.#test;
        const d = --this.#test;
        for (this.#test = 0; this.#test < 10; ++this.#test) {}
        for (this.#test = 0; this.#test < 10; this.#test++) {}

        (this.#test)++;
        (this.#test)--;
        ++(this.#test);
        --(this.#test);
        const e = (this.#test)++;
        const f = (this.#test)--;
        const g = ++(this.#test);
        const h = --(this.#test);
        for (this.#test = 0; this.#test < 10; ++(this.#test)) {}
        for (this.#test = 0; this.#test < 10; (this.#test)++) {}
    }
    test() {
        this.getInstance().#test++;
        this.getInstance().#test--;
        ++this.getInstance().#test;
        --this.getInstance().#test;
        const a = this.getInstance().#test++;
        const b = this.getInstance().#test--;
        const c = ++this.getInstance().#test;
        const d = --this.getInstance().#test;
        for (this.getInstance().#test = 0; this.getInstance().#test < 10; ++this.getInstance().#test) {}
        for (this.getInstance().#test = 0; this.getInstance().#test < 10; this.getInstance().#test++) {}

        (this.getInstance().#test)++;
        (this.getInstance().#test)--;
        ++(this.getInstance().#test);
        --(this.getInstance().#test);
        const e = (this.getInstance().#test)++;
        const f = (this.getInstance().#test)--;
        const g = ++(this.getInstance().#test);
        const h = --(this.getInstance().#test);
        for (this.getInstance().#test = 0; this.getInstance().#test < 10; ++(this.getInstance().#test)) {}
        for (this.getInstance().#test = 0; this.getInstance().#test < 10; (this.getInstance().#test)++) {}
    }
    getInstance() { return new C(); }
}


//// [privateNameFieldUnaryMutation.js]
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
var _C_test;
class C {
    constructor() {
        _C_test.set(this, 24);
        __classPrivateFieldGet(this, _C_test, "f")++;
        __classPrivateFieldGet(this, _C_test, "f")--;
        ++__classPrivateFieldGet(this, _C_test, "f");
        --__classPrivateFieldGet(this, _C_test, "f");
        const a = __classPrivateFieldGet(this, _C_test, "f")++;
        const b = __classPrivateFieldGet(this, _C_test, "f")--;
        const c = ++__classPrivateFieldGet(this, _C_test, "f");
        const d = --__classPrivateFieldGet(this, _C_test, "f");
        for (__classPrivateFieldSet(this, _C_test, 0, "f"); __classPrivateFieldGet(this, _C_test, "f") < 10; ++__classPrivateFieldGet(this, _C_test, "f")) { }
        for (__classPrivateFieldSet(this, _C_test, 0, "f"); __classPrivateFieldGet(this, _C_test, "f") < 10; __classPrivateFieldGet(this, _C_test, "f")++) { }
        (__classPrivateFieldGet(this, _C_test, "f"))++;
        (__classPrivateFieldGet(this, _C_test, "f"))--;
        ++(__classPrivateFieldGet(this, _C_test, "f"));
        --(__classPrivateFieldGet(this, _C_test, "f"));
        const e = (__classPrivateFieldGet(this, _C_test, "f"))++;
        const f = (__classPrivateFieldGet(this, _C_test, "f"))--;
        const g = ++(__classPrivateFieldGet(this, _C_test, "f"));
        const h = --(__classPrivateFieldGet(this, _C_test, "f"));
        for (__classPrivateFieldSet(this, _C_test, 0, "f"); __classPrivateFieldGet(this, _C_test, "f") < 10; ++(__classPrivateFieldGet(this, _C_test, "f"))) { }
        for (__classPrivateFieldSet(this, _C_test, 0, "f"); __classPrivateFieldGet(this, _C_test, "f") < 10; (__classPrivateFieldGet(this, _C_test, "f"))++) { }
    }
    test() {
        __classPrivateFieldGet(this.getInstance(), _C_test, "f")++;
        __classPrivateFieldGet(this.getInstance(), _C_test, "f")--;
        ++__classPrivateFieldGet(this.getInstance(), _C_test, "f");
        --__classPrivateFieldGet(this.getInstance(), _C_test, "f");
        const a = __classPrivateFieldGet(this.getInstance(), _C_test, "f")++;
        const b = __classPrivateFieldGet(this.getInstance(), _C_test, "f")--;
        const c = ++__classPrivateFieldGet(this.getInstance(), _C_test, "f");
        const d = --__classPrivateFieldGet(this.getInstance(), _C_test, "f");
        for (__classPrivateFieldSet(this.getInstance(), _C_test, 0, "f"); __classPrivateFieldGet(this.getInstance(), _C_test, "f") < 10; ++__classPrivateFieldGet(this.getInstance(), _C_test, "f")) { }
        for (__classPrivateFieldSet(this.getInstance(), _C_test, 0, "f"); __classPrivateFieldGet(this.getInstance(), _C_test, "f") < 10; __classPrivateFieldGet(this.getInstance(), _C_test, "f")++) { }
        (__classPrivateFieldGet(this.getInstance(), _C_test, "f"))++;
        (__classPrivateFieldGet(this.getInstance(), _C_test, "f"))--;
        ++(__classPrivateFieldGet(this.getInstance(), _C_test, "f"));
        --(__classPrivateFieldGet(this.getInstance(), _C_test, "f"));
        const e = (__classPrivateFieldGet(this.getInstance(), _C_test, "f"))++;
        const f = (__classPrivateFieldGet(this.getInstance(), _C_test, "f"))--;
        const g = ++(__classPrivateFieldGet(this.getInstance(), _C_test, "f"));
        const h = --(__classPrivateFieldGet(this.getInstance(), _C_test, "f"));
        for (__classPrivateFieldSet(this.getInstance(), _C_test, 0, "f"); __classPrivateFieldGet(this.getInstance(), _C_test, "f") < 10; ++(__classPrivateFieldGet(this.getInstance(), _C_test, "f"))) { }
        for (__classPrivateFieldSet(this.getInstance(), _C_test, 0, "f"); __classPrivateFieldGet(this.getInstance(), _C_test, "f") < 10; (__classPrivateFieldGet(this.getInstance(), _C_test, "f"))++) { }
    }
    getInstance() { return new C(); }
}
_C_test = new WeakMap();
