//// [tests/cases/compiler/shadowedExportsIdentifier.ts] ////

//// [main.js]
// Each `exports.a` assignment in this file is supposed to reference an `exports` local.
// If it doesn't, the `exports.a` assignment will conflict with the ESM export below.

export const foo = 1

function f1(/** @type {any} */ exports) {
    exports.a = 1
}

function f2(/** @type {any} */ { exports }) {
    exports.a = 1
}

function f3(/** @type {any} */ { x: exports }) {
    exports.a = 1
}

function f4(/** @type {any} */ [exports]) {
    exports.a = 1
}

function f10() {
    let exports = {}
    exports.a = 1
}

function f11() {
    let { exports } = { exports: { a: 0 } }
    exports.a = 1
}

function f12() {
    let  exports = {}
    while (true) {
        exports.a = 1
    }
}

function f13() {
    let exports = { a: 0 }
    function g() {
        exports.a = 1
    }
}

function f14() {
    let exports = { a: 0 }
    const g = () => {
        exports.a = 1
    }
}

function f15() {
    const exports = {}
    exports.a = 1
}

function f16() {
    var exports = {}
    exports.a = 1
}

function f17() {
    function exports() {}
    exports.a = 1
}

function f18() {
    class exports {}
    exports.a = 1
}

function f19() {
    for (const exports of [{ a: 0 }]) {
        exports.a = 1
    }
}

function f20() {
    try {
    }
    catch (/** @type {any} */ exports) {
        exports.a = 1
    }
}


//// [main.js]
// Each `exports.a` assignment in this file is supposed to reference an `exports` local.
// If it doesn't, the `exports.a` assignment will conflict with the ESM export below.
export const foo = 1;
function f1(/** @type {any} */ exports) {
    exports.a = 1;
}
function f2(/** @type {any} */ { exports }) {
    exports.a = 1;
}
function f3(/** @type {any} */ { x: exports }) {
    exports.a = 1;
}
function f4(/** @type {any} */ [exports]) {
    exports.a = 1;
}
function f10() {
    let exports = {};
    exports.a = 1;
}
function f11() {
    let { exports } = { exports: { a: 0 } };
    exports.a = 1;
}
function f12() {
    let exports = {};
    while (true) {
        exports.a = 1;
    }
}
function f13() {
    let exports = { a: 0 };
    function g() {
        exports.a = 1;
    }
}
function f14() {
    let exports = { a: 0 };
    const g = () => {
        exports.a = 1;
    };
}
function f15() {
    const exports = {};
    exports.a = 1;
}
function f16() {
    var exports = {};
    exports.a = 1;
}
function f17() {
    function exports() { }
    exports.a = 1;
}
function f18() {
    class exports {
    }
    exports.a = 1;
}
function f19() {
    for (const exports of [{ a: 0 }]) {
        exports.a = 1;
    }
}
function f20() {
    try {
    }
    catch ( /** @type {any} */exports) {
        exports.a = 1;
    }
}


//// [main.d.ts]
export declare const foo = 1;
