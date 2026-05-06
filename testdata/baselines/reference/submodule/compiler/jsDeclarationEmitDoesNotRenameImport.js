//// [tests/cases/compiler/jsDeclarationEmitDoesNotRenameImport.ts] ////

//// [Test.js]
/** @module test/Test */
class Test {}
export default Test;
//// [Test.js]
/** @module Test */
class Test {}
export default Test;
//// [index.js]
import Test from './test/Test.js'

/**
 * @typedef {Object} Options
 * @property {typeof import("./Test.js").default} [test]
 */

class X extends Test {
    /**
     * @param {Options} options
     */
    constructor(options) {
        super();
        if (options.test) {
            this.test = new options.test();
        }
    }
}

export default X;




//// [Test.d.ts]
/** @module test/Test */
class Test {
}
export default Test;
//// [Test.d.ts]
/** @module Test */
class Test {
}
export default Test;
//// [index.d.ts]
import Test from './test/Test.js';
export type Options = {
    test?: typeof import("./Test.js").default;
};
/**
 * @typedef {Object} Options
 * @property {typeof import("./Test.js").default} [test]
 */
class X extends Test {
    /**
     * @param {Options} options
     */
    constructor(options: Options);
}
export default X;


//// [DtsFileErrors]


Test.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
index.d.ts(9,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
test/Test.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== test/Test.d.ts (1 errors) ====
    /** @module test/Test */
    class Test {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    export default Test;
    
==== Test.d.ts (1 errors) ====
    /** @module Test */
    class Test {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    export default Test;
    
==== index.d.ts (1 errors) ====
    import Test from './test/Test.js';
    export type Options = {
        test?: typeof import("./Test.js").default;
    };
    /**
     * @typedef {Object} Options
     * @property {typeof import("./Test.js").default} [test]
     */
    class X extends Test {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        /**
         * @param {Options} options
         */
        constructor(options: Options);
    }
    export default X;
    