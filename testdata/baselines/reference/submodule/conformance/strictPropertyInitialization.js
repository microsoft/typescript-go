//// [tests/cases/conformance/classes/propertyMemberDeclarations/strictPropertyInitialization.ts] ////

//// [strictPropertyInitialization.ts]
// Properties with non-undefined types require initialization

class C1 {
    a: number;  // Error
    b: number | undefined;
    c: number | null;  // Error
    d?: number;
    #f: number; //Error
    #g: number | undefined;
    #h: number | null; //Error
    #i?: number;
}

// No strict initialization checks in ambient contexts

declare class C2 {
    a: number;
    b: number | undefined;
    c: number | null;
    d?: number;
    
    #f: number;
    #g: number | undefined;
    #h: number | null;
    #i?: number;
}

// No strict initialization checks for static members

class C3 {
    static a: number;
    static b: number | undefined;
    static c: number | null;
    static d?: number;
}

// Initializer satisfies strict initialization check

class C4 {
    a = 0;
    b: number = 0;
    c: string = "abc";
    #d = 0
    #e: number = 0
    #f: string= "abc"
}

// Assignment in constructor satisfies strict initialization check

class C5 {
    a: number;
    #b: number;
    constructor() {
        this.a = 0;
        this.#b = 0;
    }
}

// All code paths must contain assignment

class C6 {
    a: number;  // Error
    #b: number
    constructor(cond: boolean) {
        if (cond) {
            return;
        }
        this.a = 0;
        this.#b = 0;
    }
}

class C7 {
    a: number;
    #b: number;
    constructor(cond: boolean) {
        if (cond) {
            this.a = 1;
            this.#b = 1;
            return;
        }
        this.a = 0;
        this.#b = 1;
    }
}

// Properties with string literal names aren't checked

class C8 {
    a: number;  // Error
    "b": number;
    0: number;
}

// No strict initialization checks for abstract members

abstract class C9 {
    abstract a: number;
    abstract b: number | undefined;
    abstract c: number | null;
    abstract d?: number;
}

// Properties with non-undefined types must be assigned before they can be accessed
// within their constructor

class C10 {
    a: number;
    b: number;
    c?: number;
    #d: number;
    constructor() {
        let x = this.a;  // Error
        this.a = this.b;  // Error
        this.b = this.#d //Error
        this.b = x;
        this.#d = x;
        let y = this.c;
    }
}

// Property is considered initialized by type any even though value could be undefined

declare function someValue(): any;

class C11 {
    a: number;
    #b: number;
    constructor() {
        this.a = someValue();
        this.#b = someValue();
    }
}

const a = 'a';
const b = Symbol();

class C12 {
    [a]: number;
    [b]: number;
    ['c']: number;

    constructor() {
        this[a] = 1;
        this[b] = 1;
        this['c'] = 1;
    }
}

enum E {
    A = "A",
    B = "B"
}
class C13 {
    [E.A]: number;
    constructor() {
        this[E.A] = 1;
    }
}


//// [strictPropertyInitialization.js]
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
var _C1_f, _C1_g, _C1_h, _C1_i, _C4_d, _C4_e, _C4_f, _C5_b, _C6_b, _C7_b, _C10_d, _C11_b;
// Properties with non-undefined types require initialization
class C1 {
    constructor() {
        _C1_f.set(this, void 0); //Error
        _C1_g.set(this, void 0);
        _C1_h.set(this, void 0); //Error
        _C1_i.set(this, void 0);
    }
    a; // Error
    b;
    c; // Error
    d;
}
_C1_f = new WeakMap( // Properties with non-undefined types require initialization
// Properties with non-undefined types require initialization
), _C1_g = new WeakMap( // Properties with non-undefined types require initialization
// Properties with non-undefined types require initialization
), _C1_h = new WeakMap( // Properties with non-undefined types require initialization
// Properties with non-undefined types require initialization
), _C1_i = new WeakMap( // Properties with non-undefined types require initialization
// Properties with non-undefined types require initialization
);
// No strict initialization checks for static members
class C3 {
    static a;
    static b;
    static c;
    static d;
}
// Initializer satisfies strict initialization check
class C4 {
    constructor() {
        _C4_d.set(this, 0);
        _C4_e.set(this, 0);
        _C4_f.set(this, "abc");
    }
    a = 0;
    b = 0;
    c = "abc";
}
_C4_d = new WeakMap( // Properties with non-undefined types require initialization
// Properties with non-undefined types require initialization
), _C4_e = new WeakMap( // Properties with non-undefined types require initialization
// Properties with non-undefined types require initialization
), _C4_f = new WeakMap( // Properties with non-undefined types require initialization
// Properties with non-undefined types require initialization
);
// Assignment in constructor satisfies strict initialization check
class C5 {
    a;
    constructor() {
        _C5_b.set(this, void 0);
        this.a = 0;
        __classPrivateFieldSet(this, _C5_b, 0, "f");
    }
}
_C5_b = new WeakMap( // Properties with non-undefined types require initialization
// Properties with non-undefined types require initialization
);
// All code paths must contain assignment
class C6 {
    a; // Error
    constructor(cond) {
        _C6_b.set(this, void 0);
        if (cond) {
            return;
        }
        this.a = 0;
        __classPrivateFieldSet(this, _C6_b, 0, "f");
    }
}
_C6_b = new WeakMap( // Properties with non-undefined types require initialization
// Properties with non-undefined types require initialization
);
class C7 {
    a;
    constructor(cond) {
        _C7_b.set(this, void 0);
        if (cond) {
            this.a = 1;
            __classPrivateFieldSet(this, _C7_b, 1, "f");
            return;
        }
        this.a = 0;
        __classPrivateFieldSet(this, _C7_b, 1, "f");
    }
}
_C7_b = new WeakMap( // Properties with non-undefined types require initialization
// Properties with non-undefined types require initialization
);
// Properties with string literal names aren't checked
class C8 {
    a; // Error
    "b";
    0;
}
// No strict initialization checks for abstract members
class C9 {
    a;
    b;
    c;
    d;
}
// Properties with non-undefined types must be assigned before they can be accessed
// within their constructor
class C10 {
    a;
    b;
    c;
    constructor() {
        _C10_d.set(this, void 0);
        let x = this.a; // Error
        this.a = this.b; // Error
        this.b = __classPrivateFieldGet(this, _C10_d, "f"); //Error
        this.b = x;
        __classPrivateFieldSet(this, _C10_d, x, "f");
        let y = this.c;
    }
}
_C10_d = new WeakMap( // Properties with non-undefined types require initialization
// Properties with non-undefined types require initialization
);
class C11 {
    a;
    constructor() {
        _C11_b.set(this, void 0);
        this.a = someValue();
        __classPrivateFieldSet(this, _C11_b, someValue(), "f");
    }
}
_C11_b = new WeakMap( // Properties with non-undefined types require initialization
// Properties with non-undefined types require initialization
);
const a = 'a';
const b = Symbol();
class C12 {
    [a];
    [b];
    ['c'];
    constructor() {
        this[a] = 1;
        this[b] = 1;
        this['c'] = 1;
    }
}
var E;
(function (E) {
    E["A"] = "A";
    E["B"] = "B";
})(E || (E = {}));
class C13 {
    [E.A];
    constructor() {
        this[E.A] = 1;
    }
}


//// [strictPropertyInitialization.d.ts]
// Properties with non-undefined types require initialization
declare class C1 {
    #private;
    a: number; // Error
    b: number | undefined;
    c: number | null; // Error
    d?: number;
}
// No strict initialization checks in ambient contexts
declare class C2 {
    #private;
    a: number;
    b: number | undefined;
    c: number | null;
    d?: number;
}
// No strict initialization checks for static members
declare class C3 {
    static a: number;
    static b: number | undefined;
    static c: number | null;
    static d?: number;
}
// Initializer satisfies strict initialization check
declare class C4 {
    #private;
    a: number;
    b: number;
    c: string;
}
// Assignment in constructor satisfies strict initialization check
declare class C5 {
    #private;
    a: number;
    constructor();
}
// All code paths must contain assignment
declare class C6 {
    #private;
    a: number; // Error
    constructor(cond: boolean);
}
declare class C7 {
    #private;
    a: number;
    constructor(cond: boolean);
}
// Properties with string literal names aren't checked
declare class C8 {
    a: number; // Error
    "b": number;
    0: number;
}
// No strict initialization checks for abstract members
declare abstract class C9 {
    abstract a: number;
    abstract b: number | undefined;
    abstract c: number | null;
    abstract d?: number;
}
// Properties with non-undefined types must be assigned before they can be accessed
// within their constructor
declare class C10 {
    #private;
    a: number;
    b: number;
    c?: number;
    constructor();
}
// Property is considered initialized by type any even though value could be undefined
declare function someValue(): any;
declare class C11 {
    #private;
    a: number;
    constructor();
}
declare const a = "a";
declare const b: unique symbol;
declare class C12 {
    [a]: number;
    [b]: number;
    ['c']: number;
    constructor();
}
declare enum E {
    A = "A",
    B = "B"
}
declare class C13 {
    [E.A]: number;
    constructor();
}
