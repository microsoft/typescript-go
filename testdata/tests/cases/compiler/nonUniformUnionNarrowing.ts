// @noEmit: true

enum E { A = 1, B = 2 }

function f1(x: E | 1 | 2) {
    if (x === E.A) {
        x;  // 1 | E.A
    }
    else {
        x;  // 2 | E.B
    }
}

function f2(x: E | 1 | 2) {
    if (x === 1) {
        x;  // 1 | E.A
    }
    else {
        x;  // 2 | E.B
    }
}

namespace N1 {
    export enum E { A = 1, B = 2 }
}
namespace N2 {
    export enum E { A = 1, B = 2 }
}

function f3(x: N1.E | N2.E) {
    if (x === N1.E.A) {
        x;  // N1.E.A | N2.E.A
    }
    else {
        x;  // N1.E.B | N2.E.B
    }
}

const sym = Symbol();

function f4(x: symbol | string) {
    if (x === sym) {
        x;  // symbol (not narrowed)
    }
}
