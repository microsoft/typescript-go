//// [tests/cases/compiler/accessorDeclarationEmitJs.ts] ////

//// [a.js]
export const t1 = {
    p: 'value',
    get getter() {
        return 'value';
    },
}

export const t2 = {
    v: 'value',
    set setter(v) {},
}

export const t3 = {
    p: 'value',
    get value() {
        return 'value';
    },
    set value(v) {},
}


//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.t3 = exports.t2 = exports.t1 = void 0;
exports.t1 = {
    p: 'value',
    get getter() {
        return 'value';
    },
};
exports.t2 = {
    v: 'value',
    set setter(v) { },
};
exports.t3 = {
    p: 'value',
    get value() {
        return 'value';
    },
    set value(v) { },
};
