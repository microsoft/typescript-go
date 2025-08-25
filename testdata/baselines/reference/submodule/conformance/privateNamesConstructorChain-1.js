//// [tests/cases/conformance/classes/members/privateNames/privateNamesConstructorChain-1.ts] ////

//// [privateNamesConstructorChain-1.ts]
class Parent {
    #foo = 3;
    static #bar = 5;
    accessChildProps() {
        new Child().#foo; // OK (`#foo` was added when `Parent`'s constructor was called on `child`)
        Child.#bar;       // Error: not found
    }
}

class Child extends Parent {
    #foo = "foo";       // OK (Child's #foo does not conflict, as `Parent`'s `#foo` is not accessible)
    #bar = "bar";       // OK
}


//// [privateNamesConstructorChain-1.js]
var __classPrivateFieldGet = (this && this.__classPrivateFieldGet) || function (receiver, state, kind, f) {
    if (kind === "a" && !f) throw new TypeError("Private accessor was defined without a getter");
    if (typeof state === "function" ? receiver !== state || !f : !state.has(receiver)) throw new TypeError("Cannot read private member from an object whose class did not declare it");
    return kind === "m" ? f : kind === "a" ? f.call(receiver) : f ? f.value : state.get(receiver);
};
var _Parent_foo, _Child_foo, _Child_bar;
class Parent {
    constructor() {
        _Parent_foo.set(this, 3);
    }
    static #bar = 5;
    accessChildProps() {
        __classPrivateFieldGet(new Child(), _Parent_foo, "f"); // OK (`#foo` was added when `Parent`'s constructor was called on `child`)
        Child.#bar; // Error: not found
    }
}
_Parent_foo = new WeakMap();
class Child extends Parent {
    constructor() {
        _Child_foo.set(this, "foo"); // OK (Child's #foo does not conflict, as `Parent`'s `#foo` is not accessible)
        _Child_bar.set(this, "bar"); // OK
    }
}
_Child_foo = new WeakMap(), _Child_bar = new WeakMap();
