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
var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _o, _p, _q, _r, _s, _t, _u, _v, _w, _x, _y, _z, _0, _1, _2, _3, _4, _5, _6, _7, _8, _9, _10, _11, _12, _13, _14, _15, _16, _17, _18, _19, _20, _21, _22, _23, _24, _25, _26, _27, _28, _29, _30, _31, _32, _33, _34, _35, _36, _37;
function x(o, k) { }
let i = 0;
function foo() { return ++i + ""; }
const fieldNameA = "fieldName1";
const fieldNameB = "fieldName2";
const fieldNameC = "fieldName3";
class A {
    ["property"];
    [_a = Symbol.toStringTag];
    ["property2"] = 2;
    [_b = Symbol.iterator] = null;
    ["property3"];
    [Symbol.isConcatSpreadable];
    ["property4"] = 2;
    [Symbol.match] = null;
    [foo()];
    [_c = foo()];
    [_d = foo()] = null;
    [fieldNameA];
    [_e = fieldNameB];
    [_f = fieldNameC] = null;
}
__decorate([
    x
], A.prototype, "property", void 0);
__decorate([
    x
], A.prototype, _a, void 0);
__decorate([
    x
], A.prototype, "property2", void 0);
__decorate([
    x
], A.prototype, _b, void 0);
__decorate([
    x
], A.prototype, _c, void 0);
__decorate([
    x
], A.prototype, _d, void 0);
__decorate([
    x
], A.prototype, _e, void 0);
__decorate([
    x
], A.prototype, _f, void 0);
void class B {
    ["property"];
    [_g = Symbol.toStringTag];
    ["property2"] = 2;
    [_h = Symbol.iterator] = null;
    ["property3"];
    [Symbol.isConcatSpreadable];
    ["property4"] = 2;
    [Symbol.match] = null;
    [foo()];
    [_j = foo()];
    [_k = foo()] = null;
    [fieldNameA];
    [_l = fieldNameB];
    [_m = fieldNameC] = null;
};
class C {
    ["property"];
    [_o = Symbol.toStringTag];
    ["property2"] = 2;
    [_p = Symbol.iterator] = null;
    ["property3"];
    [Symbol.isConcatSpreadable];
    ["property4"] = 2;
    [Symbol.match] = null;
    [foo()];
    [_q = foo()];
    [_r = foo()] = null;
    [fieldNameA];
    [_s = fieldNameB];
    [_t = fieldNameC] = null;
    ["some" + "method"]() { }
}
__decorate([
    x
], C.prototype, "property", void 0);
__decorate([
    x
], C.prototype, _o, void 0);
__decorate([
    x
], C.prototype, "property2", void 0);
__decorate([
    x
], C.prototype, _p, void 0);
__decorate([
    x
], C.prototype, _q, void 0);
__decorate([
    x
], C.prototype, _r, void 0);
__decorate([
    x
], C.prototype, _s, void 0);
__decorate([
    x
], C.prototype, _t, void 0);
void class D {
    ["property"];
    [_u = Symbol.toStringTag];
    ["property2"] = 2;
    [_v = Symbol.iterator] = null;
    ["property3"];
    [Symbol.isConcatSpreadable];
    ["property4"] = 2;
    [Symbol.match] = null;
    [foo()];
    [_w = foo()];
    [_x = foo()] = null;
    [fieldNameA];
    [_y = fieldNameB];
    [_z = fieldNameC] = null;
    ["some" + "method"]() { }
};
class E {
    ["property"];
    [_0 = Symbol.toStringTag];
    ["property2"] = 2;
    [_1 = Symbol.iterator] = null;
    ["property3"];
    [Symbol.isConcatSpreadable];
    ["property4"] = 2;
    [Symbol.match] = null;
    [foo()];
    [_2 = foo()];
    [_3 = foo()] = null;
    ["some" + "method"]() { }
    [fieldNameA];
    [_4 = fieldNameB];
    [_5 = fieldNameC] = null;
}
__decorate([
    x
], E.prototype, "property", void 0);
__decorate([
    x
], E.prototype, _0, void 0);
__decorate([
    x
], E.prototype, "property2", void 0);
__decorate([
    x
], E.prototype, _1, void 0);
__decorate([
    x
], E.prototype, _2, void 0);
__decorate([
    x
], E.prototype, _3, void 0);
__decorate([
    x
], E.prototype, _4, void 0);
__decorate([
    x
], E.prototype, _5, void 0);
void class F {
    ["property"];
    [_6 = Symbol.toStringTag];
    ["property2"] = 2;
    [_7 = Symbol.iterator] = null;
    ["property3"];
    [Symbol.isConcatSpreadable];
    ["property4"] = 2;
    [Symbol.match] = null;
    [foo()];
    [_8 = foo()];
    [_9 = foo()] = null;
    ["some" + "method"]() { }
    [fieldNameA];
    [_10 = fieldNameB];
    [_11 = fieldNameC] = null;
};
class G {
    ["property"];
    [_12 = Symbol.toStringTag];
    ["property2"] = 2;
    [_13 = Symbol.iterator] = null;
    ["property3"];
    [Symbol.isConcatSpreadable];
    ["property4"] = 2;
    [Symbol.match] = null;
    [foo()];
    [_14 = foo()];
    [_15 = foo()] = null;
    ["some" + "method"]() { }
    [fieldNameA];
    [_16 = fieldNameB];
    ["some" + "method2"]() { }
    [_17 = fieldNameC] = null;
}
__decorate([
    x
], G.prototype, "property", void 0);
__decorate([
    x
], G.prototype, _12, void 0);
__decorate([
    x
], G.prototype, "property2", void 0);
__decorate([
    x
], G.prototype, _13, void 0);
__decorate([
    x
], G.prototype, _14, void 0);
__decorate([
    x
], G.prototype, _15, void 0);
__decorate([
    x
], G.prototype, _16, void 0);
__decorate([
    x
], G.prototype, _17, void 0);
void class H {
    ["property"];
    [_18 = Symbol.toStringTag];
    ["property2"] = 2;
    [_19 = Symbol.iterator] = null;
    ["property3"];
    [Symbol.isConcatSpreadable];
    ["property4"] = 2;
    [Symbol.match] = null;
    [foo()];
    [_20 = foo()];
    [_21 = foo()] = null;
    ["some" + "method"]() { }
    [fieldNameA];
    [_22 = fieldNameB];
    ["some" + "method2"]() { }
    [_23 = fieldNameC] = null;
};
class I {
    ["property"];
    [_24 = Symbol.toStringTag];
    ["property2"] = 2;
    [_25 = Symbol.iterator] = null;
    ["property3"];
    [Symbol.isConcatSpreadable];
    ["property4"] = 2;
    [Symbol.match] = null;
    [foo()];
    [_26 = foo()];
    [_27 = foo()] = null;
    [_28 = "some" + "method"]() { }
    [fieldNameA];
    [_29 = fieldNameB];
    ["some" + "method2"]() { }
    [_30 = fieldNameC] = null;
}
__decorate([
    x
], I.prototype, "property", void 0);
__decorate([
    x
], I.prototype, _24, void 0);
__decorate([
    x
], I.prototype, "property2", void 0);
__decorate([
    x
], I.prototype, _25, void 0);
__decorate([
    x
], I.prototype, _26, void 0);
__decorate([
    x
], I.prototype, _27, void 0);
__decorate([
    x
], I.prototype, _28, null);
__decorate([
    x
], I.prototype, _29, void 0);
__decorate([
    x
], I.prototype, _30, void 0);
void class J {
    ["property"];
    [_31 = Symbol.toStringTag];
    ["property2"] = 2;
    [_32 = Symbol.iterator] = null;
    ["property3"];
    [Symbol.isConcatSpreadable];
    ["property4"] = 2;
    [Symbol.match] = null;
    [foo()];
    [_33 = foo()];
    [_34 = foo()] = null;
    [_35 = "some" + "method"]() { }
    [fieldNameA];
    [_36 = fieldNameB];
    ["some" + "method2"]() { }
    [_37 = fieldNameC] = null;
};
