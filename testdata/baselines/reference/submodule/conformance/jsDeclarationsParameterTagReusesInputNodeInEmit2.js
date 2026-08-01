//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsParameterTagReusesInputNodeInEmit2.ts] ////

//// [base.js]
class Base {
    constructor() {}
}

const BaseFactory = () => {
    return new Base();
};

BaseFactory.Base = Base;

module.exports = BaseFactory;

//// [file.js]
/** @typedef {typeof import('./base')} BaseFactory */

/**
 *
 * @param {InstanceType<BaseFactory["Base"]>} base
 * @returns {InstanceType<BaseFactory["Base"]>}
 */
const test = (base) => {
    return base;
};


//// [base.js]
"use strict";
class Base {
    constructor() { }
}
const BaseFactory = () => {
    return new Base();
};
BaseFactory.Base = Base;
module.exports = BaseFactory;
//// [file.js]
"use strict";
/** @typedef {typeof import('./base')} BaseFactory */
/**
 *
 * @param {InstanceType<BaseFactory["Base"]>} base
 * @returns {InstanceType<BaseFactory["Base"]>}
 */
const test = (base) => {
    return base;
};


//// [base.d.ts]
export = _exports;
declare function _exports(): Base;
declare class Base {
    constructor();
}
declare function BaseFactory(): Base;
declare namespace BaseFactory {
    export { Base };
}
//// [file.d.ts]
/** @typedef {typeof import('./base')} BaseFactory */
type BaseFactory = typeof import('./base');
/**
 *
 * @param {InstanceType<BaseFactory["Base"]>} base
 * @returns {InstanceType<BaseFactory["Base"]>}
 */
declare const test: (base: InstanceType<BaseFactory["Base"]>) => InstanceType<BaseFactory["Base"]>;


//// [DtsFileErrors]


out/file.d.ts(8,53): error TS2339: Property 'Base' does not exist on type '() => Base'.
out/file.d.ts(8,91): error TS2339: Property 'Base' does not exist on type '() => Base'.


==== out/base.d.ts (0 errors) ====
    export = _exports;
    declare function _exports(): Base;
    declare class Base {
        constructor();
    }
    declare function BaseFactory(): Base;
    declare namespace BaseFactory {
        export { Base };
    }
    
==== out/file.d.ts (2 errors) ====
    /** @typedef {typeof import('./base')} BaseFactory */
    type BaseFactory = typeof import('./base');
    /**
     *
     * @param {InstanceType<BaseFactory["Base"]>} base
     * @returns {InstanceType<BaseFactory["Base"]>}
     */
    declare const test: (base: InstanceType<BaseFactory["Base"]>) => InstanceType<BaseFactory["Base"]>;
                                                        ~~~~~~
!!! error TS2339: Property 'Base' does not exist on type '() => Base'.
                                                                                              ~~~~~~
!!! error TS2339: Property 'Base' does not exist on type '() => Base'.
    