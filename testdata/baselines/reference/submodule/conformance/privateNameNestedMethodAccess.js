//// [tests/cases/conformance/classes/members/privateNames/privateNameNestedMethodAccess.ts] ////

//// [privateNameNestedMethodAccess.ts]
class C {
    #foo = 42;
    #bar() { new C().#baz; }
    get #baz() { return 42; }

    m() {
        return class D {
            #bar() {}
            constructor() {
                new C().#foo;
                new C().#bar; // Error
                new C().#baz;
                new D().#bar;
            }

            n(x: any) {
                x.#foo;
                x.#bar;
                x.#unknown; // Error
            }
        }
    }
}


//// [privateNameNestedMethodAccess.js]
var __classPrivateFieldGet = (this && this.__classPrivateFieldGet) || function (receiver, state, kind, f) {
    if (kind === "a" && !f) throw new TypeError("Private accessor was defined without a getter");
    if (typeof state === "function" ? receiver !== state || !f : !state.has(receiver)) throw new TypeError("Cannot read private member from an object whose class did not declare it");
    return kind === "m" ? f : kind === "a" ? f.call(receiver) : f ? f.value : state.get(receiver);
};
var _C_foo;
class C {
    constructor() {
        _C_foo.set(this, 42);
    }
    #bar() { new C().#baz; }
    get #baz() { return 42; }
    m() {
        return class D {
            #bar() { }
            constructor() {
                __classPrivateFieldGet(new C(), _C_foo, "f");
                new C().#bar; // Error
                new C().#baz;
                new D().#bar;
            }
            n(x) {
                __classPrivateFieldGet(x, _C_foo, "f");
                x.#bar;
                x.#unknown; // Error
            }
        };
    }
}
_C_foo = new WeakMap();
