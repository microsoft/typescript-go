//// [tests/cases/conformance/es7/exponentiationOperator/emitCompoundExponentiationAssignmentWithIndexingOnLHS2.ts] ////

//// [emitCompoundExponentiationAssignmentWithIndexingOnLHS2.ts]
var globalCounter = 0;
function foo() {
    globalCounter += 1;
    return { 0: 2 };
}
foo()[0] **= foo()[0];
var result_foo1 = foo()[0] **= foo()[0];
foo()[0] **= foo()[0] **= 2;
var result_foo2 = foo()[0] **= foo()[0] **= 2;
foo()[0] **= foo()[0] ** 2;
var result_foo3 = foo()[0] **= foo()[0] ** 2;

//// [emitCompoundExponentiationAssignmentWithIndexingOnLHS2.js]
var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k, _l, _m, _n, _o, _p, _q;
var globalCounter = 0;
function foo() {
    globalCounter += 1;
    return { 0: 2 };
}
(_a = foo())[_b = 0] = Math.pow(_a[_b], foo()[0]);
var result_foo1 = (_c = foo())[_d = 0] = Math.pow(_c[_d], foo()[0]);
(_g = foo())[_h = 0] = Math.pow(_g[_h], (_e = foo())[_f = 0] = Math.pow(_e[_f], 2));
var result_foo2 = (_l = foo())[_m = 0] = Math.pow(_l[_m], (_j = foo())[_k = 0] = Math.pow(_j[_k], 2));
(_n = foo())[_o = 0] = Math.pow(_n[_o], Math.pow(foo()[0], 2));
var result_foo3 = (_p = foo())[_q = 0] = Math.pow(_p[_q], Math.pow(foo()[0], 2));
