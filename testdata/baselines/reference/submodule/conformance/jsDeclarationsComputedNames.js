//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsComputedNames.ts] ////

//// [index.js]
const TopLevelSym = Symbol();
const InnerSym = Symbol();
module.exports = {
    [TopLevelSym](x = 12) {
        return x;
    },
    items: {
        [InnerSym]: (arg = {x: 12}) => arg.x
    }
}

//// [index2.js]
const TopLevelSym = Symbol();
const InnerSym = Symbol();

export class MyClass {
    static [TopLevelSym] = 12;
    [InnerSym] = "ok";
    /**
     * @param {typeof TopLevelSym | typeof InnerSym} _p
     */
    constructor(_p = InnerSym) {
        // switch on _p
    }
}


//// [index.js]
const TopLevelSym = Symbol();
const InnerSym = Symbol();
module.exports = {
    [TopLevelSym](x = 12) {
        return x;
    },
    items: {
        [InnerSym]: (arg = { x: 12 }) => arg.x
    }
};
//// [index2.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.MyClass = void 0;
const TopLevelSym = Symbol();
const InnerSym = Symbol();
class MyClass {
    static [TopLevelSym] = 12;
    [InnerSym] = "ok";
    constructor(_p = InnerSym) {
    }
}
exports.MyClass = MyClass;
