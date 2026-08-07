//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsTypedefPropertyComments.ts] ////

//// [lib.js]
/**
 * @typedef {Object} Foo
 * @property {boolean} bool Whether `.bool` is true or not
 */
export class C {
    /** @returns {Foo} */
    getFoo() { return { bool: false }; }
}

//// [main.js]
import { C } from './lib.js';

export class Main {
    constructor() {
        this.c = new C();
    }

    getFoo() { return { ...this.c.getFoo() }; }
}




//// [lib.d.ts]
export type Foo = {
    /**
     * Whether `.bool` is true or not
     */
    bool: boolean;
};
/**
 * @typedef {Object} Foo
 * @property {boolean} bool Whether `.bool` is true or not
 */
export declare class C {
    /** @returns {Foo} */
    getFoo(): Foo;
}
//// [main.d.ts]
import { C } from './lib.js';
export declare class Main {
    c: C;
    constructor();
    getFoo(): {
        /**
         * Whether `.bool` is true or not
         */
        bool: boolean;
    };
}
