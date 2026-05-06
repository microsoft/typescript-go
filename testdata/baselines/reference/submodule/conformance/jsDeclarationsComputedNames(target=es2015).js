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
"use strict";
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
var _a, _b;
Object.defineProperty(exports, "__esModule", { value: true });
exports.MyClass = void 0;
const TopLevelSym = Symbol();
const InnerSym = Symbol();
class MyClass {
    /**
     * @param {typeof TopLevelSym | typeof InnerSym} _p
     */
    constructor(_p = InnerSym) {
        this[_b] = "ok";
        // switch on _p
    }
}
exports.MyClass = MyClass;
_a = TopLevelSym, _b = InnerSym;
MyClass[_a] = 12;


//// [index.d.ts]
const TopLevelSym: unique symbol;
const InnerSym: unique symbol;
const _default: {
    [TopLevelSym](x?: number): number;
    items: {
        [InnerSym]: (arg?: {
            x: number;
        }) => number;
    };
};
export = _default;
//// [index2.d.ts]
const TopLevelSym: unique symbol;
const InnerSym: unique symbol;
export class MyClass {
    static [TopLevelSym]: number;
    [InnerSym]: string;
    /**
     * @param {typeof TopLevelSym | typeof InnerSym} _p
     */
    constructor(_p?: typeof TopLevelSym | typeof InnerSym);
}
export {};


//// [DtsFileErrors]


out/index.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
out/index2.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/index.d.ts (1 errors) ====
    const TopLevelSym: unique symbol;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    const InnerSym: unique symbol;
    const _default: {
        [TopLevelSym](x?: number): number;
        items: {
            [InnerSym]: (arg?: {
                x: number;
            }) => number;
        };
    };
    export = _default;
    
==== out/index2.d.ts (1 errors) ====
    const TopLevelSym: unique symbol;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    const InnerSym: unique symbol;
    export class MyClass {
        static [TopLevelSym]: number;
        [InnerSym]: string;
        /**
         * @param {typeof TopLevelSym | typeof InnerSym} _p
         */
        constructor(_p?: typeof TopLevelSym | typeof InnerSym);
    }
    export {};
    