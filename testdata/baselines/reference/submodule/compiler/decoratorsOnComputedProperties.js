//// [tests/cases/compiler/decoratorsOnComputedProperties.ts] ////

//// [decoratorsOnComputedProperties.ts]
function x(o: object, k: PropertyKey) { }
let i = 0;
function foo(): string { return ++i + ""; }

const fieldNameA: string = "fieldName1";
const fieldNameB: string = "fieldName2";
const fieldNameC: string = "fieldName3";

class A {
    @x ["property"]: any;
    @x [Symbol.toStringTag]: any;
    @x ["property2"]: any = 2;
    @x [Symbol.iterator]: any = null;
    ["property3"]: any;
    [Symbol.isConcatSpreadable]: any;
    ["property4"]: any = 2;
    [Symbol.match]: any = null;
    [foo()]: any;
    @x [foo()]: any;
    @x [foo()]: any = null;
    [fieldNameA]: any;
    @x [fieldNameB]: any;
    @x [fieldNameC]: any = null;
}

void class B {
    @x ["property"]: any;
    @x [Symbol.toStringTag]: any;
    @x ["property2"]: any = 2;
    @x [Symbol.iterator]: any = null;
    ["property3"]: any;
    [Symbol.isConcatSpreadable]: any;
    ["property4"]: any = 2;
    [Symbol.match]: any = null;
    [foo()]: any;
    @x [foo()]: any;
    @x [foo()]: any = null;
    [fieldNameA]: any;
    @x [fieldNameB]: any;
    @x [fieldNameC]: any = null;
};

class C {
    @x ["property"]: any;
    @x [Symbol.toStringTag]: any;
    @x ["property2"]: any = 2;
    @x [Symbol.iterator]: any = null;
    ["property3"]: any;
    [Symbol.isConcatSpreadable]: any;
    ["property4"]: any = 2;
    [Symbol.match]: any = null;
    [foo()]: any;
    @x [foo()]: any;
    @x [foo()]: any = null;
    [fieldNameA]: any;
    @x [fieldNameB]: any;
    @x [fieldNameC]: any = null;
    ["some" + "method"]() {}
}

void class D {
    @x ["property"]: any;
    @x [Symbol.toStringTag]: any;
    @x ["property2"]: any = 2;
    @x [Symbol.iterator]: any = null;
    ["property3"]: any;
    [Symbol.isConcatSpreadable]: any;
    ["property4"]: any = 2;
    [Symbol.match]: any = null;
    [foo()]: any;
    @x [foo()]: any;
    @x [foo()]: any = null;
    [fieldNameA]: any;
    @x [fieldNameB]: any;
    @x [fieldNameC]: any = null;
    ["some" + "method"]() {}
};

class E {
    @x ["property"]: any;
    @x [Symbol.toStringTag]: any;
    @x ["property2"]: any = 2;
    @x [Symbol.iterator]: any = null;
    ["property3"]: any;
    [Symbol.isConcatSpreadable]: any;
    ["property4"]: any = 2;
    [Symbol.match]: any = null;
    [foo()]: any;
    @x [foo()]: any;
    @x [foo()]: any = null;
    ["some" + "method"]() {}
    [fieldNameA]: any;
    @x [fieldNameB]: any;
    @x [fieldNameC]: any = null;
}

void class F {
    @x ["property"]: any;
    @x [Symbol.toStringTag]: any;
    @x ["property2"]: any = 2;
    @x [Symbol.iterator]: any = null;
    ["property3"]: any;
    [Symbol.isConcatSpreadable]: any;
    ["property4"]: any = 2;
    [Symbol.match]: any = null;
    [foo()]: any;
    @x [foo()]: any;
    @x [foo()]: any = null;
    ["some" + "method"]() {}
    [fieldNameA]: any;
    @x [fieldNameB]: any;
    @x [fieldNameC]: any = null;
};

class G {
    @x ["property"]: any;
    @x [Symbol.toStringTag]: any;
    @x ["property2"]: any = 2;
    @x [Symbol.iterator]: any = null;
    ["property3"]: any;
    [Symbol.isConcatSpreadable]: any;
    ["property4"]: any = 2;
    [Symbol.match]: any = null;
    [foo()]: any;
    @x [foo()]: any;
    @x [foo()]: any = null;
    ["some" + "method"]() {}
    [fieldNameA]: any;
    @x [fieldNameB]: any;
    ["some" + "method2"]() {}
    @x [fieldNameC]: any = null;
}

void class H {
    @x ["property"]: any;
    @x [Symbol.toStringTag]: any;
    @x ["property2"]: any = 2;
    @x [Symbol.iterator]: any = null;
    ["property3"]: any;
    [Symbol.isConcatSpreadable]: any;
    ["property4"]: any = 2;
    [Symbol.match]: any = null;
    [foo()]: any;
    @x [foo()]: any;
    @x [foo()]: any = null;
    ["some" + "method"]() {}
    [fieldNameA]: any;
    @x [fieldNameB]: any;
    ["some" + "method2"]() {}
    @x [fieldNameC]: any = null;
};

class I {
    @x ["property"]: any;
    @x [Symbol.toStringTag]: any;
    @x ["property2"]: any = 2;
    @x [Symbol.iterator]: any = null;
    ["property3"]: any;
    [Symbol.isConcatSpreadable]: any;
    ["property4"]: any = 2;
    [Symbol.match]: any = null;
    [foo()]: any;
    @x [foo()]: any;
    @x [foo()]: any = null;
    @x ["some" + "method"]() {}
    [fieldNameA]: any;
    @x [fieldNameB]: any;
    ["some" + "method2"]() {}
    @x [fieldNameC]: any = null;
}

void class J {
    @x ["property"]: any;
    @x [Symbol.toStringTag]: any;
    @x ["property2"]: any = 2;
    @x [Symbol.iterator]: any = null;
    ["property3"]: any;
    [Symbol.isConcatSpreadable]: any;
    ["property4"]: any = 2;
    [Symbol.match]: any = null;
    [foo()]: any;
    @x [foo()]: any;
    @x [foo()]: any = null;
    @x ["some" + "method"]() {}
    [fieldNameA]: any;
    @x [fieldNameB]: any;
    ["some" + "method2"]() {}
    @x [fieldNameC]: any = null;
};

//// [decoratorsOnComputedProperties.js]
"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s, _t, _u, _v, _w, _x, _y, _z, _0, _1, _2, _3, _4, _5, _6, _7, _8, _9, _10, _11, _12, _13, _14, _15, _16, _17, _18, _19;
function x(o, k) { }
let i = 0;
function foo() { return ++i + ""; }
const fieldNameA = "fieldName1";
const fieldNameB = "fieldName2";
const fieldNameC = "fieldName3";
class A {
    constructor() {
        this["property2"] = 2;
        this[_a] = null;
        this["property4"] = 2;
        this[_b] = null;
        this[_c] = null;
        this[_d] = null;
    }
}
Symbol.toStringTag, _a = Symbol.iterator, Symbol.isConcatSpreadable, _b = Symbol.match, foo(), foo(), _c = foo(), _d = fieldNameC;
__decorate([
    x
], A.prototype, "property", void 0);
__decorate([
    x
], A.prototype, _20, void 0);
__decorate([
    x
], A.prototype, "property2", void 0);
__decorate([
    x
], A.prototype, _a, void 0);
__decorate([
    x
], A.prototype, _21, void 0);
__decorate([
    x
], A.prototype, _c, void 0);
__decorate([
    x
], A.prototype, _22, void 0);
__decorate([
    x
], A.prototype, _d, void 0);
void (_j = class B {
        constructor() {
            this["property2"] = 2;
            this[_e] = null;
            this["property4"] = 2;
            this[_f] = null;
            this[_g] = null;
            this[_h] = null;
        }
    },
    Symbol.toStringTag,
    _e = Symbol.iterator,
    Symbol.isConcatSpreadable,
    _f = Symbol.match,
    foo(),
    foo(),
    _g = foo(),
    _h = fieldNameC,
    _j);
class C {
    constructor() {
        this["property2"] = 2;
        this[_k] = null;
        this["property4"] = 2;
        this[_l] = null;
        this[_m] = null;
        this[_o] = null;
    }
    [(Symbol.toStringTag, _k = Symbol.iterator, Symbol.isConcatSpreadable, _l = Symbol.match, foo(), foo(), _m = foo(), _o = fieldNameC, "some" + "method")]() { }
}
__decorate([
    x
], C.prototype, "property", void 0);
__decorate([
    x
], C.prototype, _23, void 0);
__decorate([
    x
], C.prototype, "property2", void 0);
__decorate([
    x
], C.prototype, _k, void 0);
__decorate([
    x
], C.prototype, _24, void 0);
__decorate([
    x
], C.prototype, _m, void 0);
__decorate([
    x
], C.prototype, _25, void 0);
__decorate([
    x
], C.prototype, _o, void 0);
void class D {
    constructor() {
        this["property2"] = 2;
        this[_p] = null;
        this["property4"] = 2;
        this[_q] = null;
        this[_r] = null;
        this[_s] = null;
    }
    [(Symbol.toStringTag, _p = Symbol.iterator, Symbol.isConcatSpreadable, _q = Symbol.match, foo(), foo(), _r = foo(), _s = fieldNameC, "some" + "method")]() { }
};
class E {
    constructor() {
        this["property2"] = 2;
        this[_t] = null;
        this["property4"] = 2;
        this[_u] = null;
        this[_v] = null;
        this[_w] = null;
    }
    [(Symbol.toStringTag, _t = Symbol.iterator, Symbol.isConcatSpreadable, _u = Symbol.match, foo(), foo(), _v = foo(), "some" + "method")]() { }
}
_w = fieldNameC;
__decorate([
    x
], E.prototype, "property", void 0);
__decorate([
    x
], E.prototype, _26, void 0);
__decorate([
    x
], E.prototype, "property2", void 0);
__decorate([
    x
], E.prototype, _t, void 0);
__decorate([
    x
], E.prototype, _27, void 0);
__decorate([
    x
], E.prototype, _v, void 0);
__decorate([
    x
], E.prototype, _28, void 0);
__decorate([
    x
], E.prototype, _w, void 0);
void (_1 = class F {
        constructor() {
            this["property2"] = 2;
            this[_x] = null;
            this["property4"] = 2;
            this[_y] = null;
            this[_z] = null;
            this[_0] = null;
        }
        [(Symbol.toStringTag, _x = Symbol.iterator, Symbol.isConcatSpreadable, _y = Symbol.match, foo(), foo(), _z = foo(), "some" + "method")]() { }
    },
    _0 = fieldNameC,
    _1);
class G {
    constructor() {
        this["property2"] = 2;
        this[_2] = null;
        this["property4"] = 2;
        this[_3] = null;
        this[_4] = null;
        this[_5] = null;
    }
    [(Symbol.toStringTag, _2 = Symbol.iterator, Symbol.isConcatSpreadable, _3 = Symbol.match, foo(), foo(), _4 = foo(), "some" + "method")]() { }
    ["some" + "method2"]() { }
}
_5 = fieldNameC;
__decorate([
    x
], G.prototype, "property", void 0);
__decorate([
    x
], G.prototype, _29, void 0);
__decorate([
    x
], G.prototype, "property2", void 0);
__decorate([
    x
], G.prototype, _2, void 0);
__decorate([
    x
], G.prototype, _30, void 0);
__decorate([
    x
], G.prototype, _4, void 0);
__decorate([
    x
], G.prototype, _31, void 0);
__decorate([
    x
], G.prototype, _5, void 0);
void (_10 = class H {
        constructor() {
            this["property2"] = 2;
            this[_6] = null;
            this["property4"] = 2;
            this[_7] = null;
            this[_8] = null;
            this[_9] = null;
        }
        [(Symbol.toStringTag, _6 = Symbol.iterator, Symbol.isConcatSpreadable, _7 = Symbol.match, foo(), foo(), _8 = foo(), "some" + "method")]() { }
        ["some" + "method2"]() { }
    },
    _9 = fieldNameC,
    _10);
class I {
    constructor() {
        this["property2"] = 2;
        this[_11] = null;
        this["property4"] = 2;
        this[_12] = null;
        this[_13] = null;
        this[_14] = null;
    }
    [(Symbol.toStringTag, _11 = Symbol.iterator, Symbol.isConcatSpreadable, _12 = Symbol.match, foo(), foo(), _13 = foo(), "some" + "method")]() { }
    ["some" + "method2"]() { }
}
_14 = fieldNameC;
__decorate([
    x
], I.prototype, "property", void 0);
__decorate([
    x
], I.prototype, _32, void 0);
__decorate([
    x
], I.prototype, "property2", void 0);
__decorate([
    x
], I.prototype, _11, void 0);
__decorate([
    x
], I.prototype, _33, void 0);
__decorate([
    x
], I.prototype, _13, void 0);
__decorate([
    x
], I.prototype, _34, null);
__decorate([
    x
], I.prototype, _35, void 0);
__decorate([
    x
], I.prototype, _14, void 0);
void (_19 = class J {
        constructor() {
            this["property2"] = 2;
            this[_15] = null;
            this["property4"] = 2;
            this[_16] = null;
            this[_17] = null;
            this[_18] = null;
        }
        [(Symbol.toStringTag, _15 = Symbol.iterator, Symbol.isConcatSpreadable, _16 = Symbol.match, foo(), foo(), _17 = foo(), "some" + "method")]() { }
        ["some" + "method2"]() { }
    },
    _18 = fieldNameC,
    _19);
