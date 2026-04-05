//// [tests/cases/compiler/shadowedModuleIdentifier.ts] ////

//// [main.js]
// Each `module.exports` assignment in this file is supposed to reference a `module` local.
// If it doesn't, the `module.exports` assignment will conflict with the ESM export below.

export const foo = 1

function f1(/** @type {any} */ module) {
    module.exports = 1
}

function f2(/** @type {any} */ { module }) {
    module.exports = 1
}

function f3(/** @type {any} */ { x: module }) {
    module.exports = 1
}

function f4(/** @type {any} */ [module]) {
    module.exports = 1
}

function f5(/** @type {any} */ module) {
    module.exports.a = 1
}

function f6(/** @type {any} */ { module }) {
    module.exports.a = 1
}

function f7(/** @type {any} */ { x: module }) {
    module.exports.a = 1
}

function f8(/** @type {any} */ [module]) {
    module.exports.a = 1
}

function f10() {
    let module = {}
    module.exports = 1
}

function f11() {
    let { module } = { module: { exports: 0 } }
    module.exports = 1
}

function f12() {
    let  module = {}
    while (true) {
        module.exports = 1
    }
}

function f13() {
    let module = { exports: 0 }
    function g() {
        module.exports = 1
    }
}

function f14() {
    let module = { exports: 0 }
    const g = () => {
        module.exports = 1
    }
}

function f15() {
    const module = {}
    module.exports = 1
}

function f16() {
    var module = {}
    module.exports = 1
}

function f17() {
    function module() {}
    module.exports = 1
}

function f18() {
    class module {}
    module.exports = 1
}

function f19() {
    for (const module of [{ exports: 0 }]) {
        module.exports = 1
    }
}

function f20() {
    try {
    }
    catch (/** @type {any} */ module) {
        module.exports = 1
    }
}

function f21() {
    /** @param {any} module */
    const g = module => {
        module.exports = 1
    }
}

function f22() {
    const g = (/** @type {any} */ module) => {
        module.exports = 1
    }
}


//// [main.js]
// Each `module.exports` assignment in this file is supposed to reference a `module` local.
// If it doesn't, the `module.exports` assignment will conflict with the ESM export below.
export const foo = 1;
function f1(/** @type {any} */ module) {
    module.exports = 1;
}
function f2(/** @type {any} */ { module }) {
    module.exports = 1;
}
function f3(/** @type {any} */ { x: module }) {
    module.exports = 1;
}
function f4(/** @type {any} */ [module]) {
    module.exports = 1;
}
function f5(/** @type {any} */ module) {
    module.exports.a = 1;
}
function f6(/** @type {any} */ { module }) {
    module.exports.a = 1;
}
function f7(/** @type {any} */ { x: module }) {
    module.exports.a = 1;
}
function f8(/** @type {any} */ [module]) {
    module.exports.a = 1;
}
function f10() {
    let module = {};
    module.exports = 1;
}
function f11() {
    let { module } = { module: { exports: 0 } };
    module.exports = 1;
}
function f12() {
    let module = {};
    while (true) {
        module.exports = 1;
    }
}
function f13() {
    let module = { exports: 0 };
    function g() {
        module.exports = 1;
    }
}
function f14() {
    let module = { exports: 0 };
    const g = () => {
        module.exports = 1;
    };
}
function f15() {
    const module = {};
    module.exports = 1;
}
function f16() {
    var module = {};
    module.exports = 1;
}
function f17() {
    function module() { }
    module.exports = 1;
}
function f18() {
    class module {
    }
    module.exports = 1;
}
function f19() {
    for (const module of [{ exports: 0 }]) {
        module.exports = 1;
    }
}
function f20() {
    try {
    }
    catch ( /** @type {any} */module) {
        module.exports = 1;
    }
}
function f21() {
    /** @param {any} module */
    const g = module => {
        module.exports = 1;
    };
}
function f22() {
    const g = (/** @type {any} */ module) => {
        module.exports = 1;
    };
}


//// [main.d.ts]
export declare const foo = 1;
