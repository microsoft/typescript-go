//// [tests/cases/conformance/decorators/class/method/decoratorOnClassMethod19.ts] ////

//// [decoratorOnClassMethod19.ts]
// https://github.com/microsoft/TypeScript/issues/48515
declare var decorator: any;

class C1 {
    #x

    @decorator((x: C1) => x.#x)
    y() {}
}

class C2 {
    #x

    y(@decorator((x: C2) => x.#x) p) {}
}


//// [decoratorOnClassMethod19.js]
var __classPrivateFieldGet = (this && this.__classPrivateFieldGet) || function (receiver, state, kind, f) {
    if (kind === "a" && !f) throw new TypeError("Private accessor was defined without a getter");
    if (typeof state === "function" ? receiver !== state || !f : !state.has(receiver)) throw new TypeError("Cannot read private member from an object whose class did not declare it");
    return kind === "m" ? f : kind === "a" ? f.call(receiver) : f ? f.value : state.get(receiver);
};
var _C1_x, _C2_x;
class C1 {
    constructor() {
        _C1_x.set(this, void 0);
    }
    @decorator((x) => __classPrivateFieldGet(x, _C1_x, "f"))
    y() { }
}
_C1_x = new WeakMap( // https://github.com/microsoft/TypeScript/issues/48515
// https://github.com/microsoft/TypeScript/issues/48515
);
class C2 {
    constructor() {
        _C2_x.set(this, void 0);
    }
    y(p) { }
}
_C2_x = new WeakMap( // https://github.com/microsoft/TypeScript/issues/48515
// https://github.com/microsoft/TypeScript/issues/48515
);
