//// [tests/cases/compiler/declarationEmitRetainsJsdocyComments.ts] ////

//// [declarationEmitRetainsJsdocyComments.ts]
/**
 * comment1
 * @param p 
 */
export const foo = (p: string) => {
    return {
        /**
         * comment2
         * @param s 
         */
        bar: (s: number) => {},
        /**
         * comment3
         * @param s 
         */
        bar2(s: number) {},
    }
}

export class Foo {
    /**
     * comment4
     * @param s  
     */
    bar(s: number) {
    }
}

export let {
    /**
    * comment5
    */
    someMethod
} = null as any;

declare global {
    interface ExtFunc {
        /**
        * comment6
        */
        someMethod(collection: any[]): boolean;
    }
}


//// [declarationEmitRetainsJsdocyComments.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.someMethod = exports.Foo = exports.foo = void 0;
/**
 * comment1
 * @param p
 */
const foo = (p) => {
    return {
        /**
         * comment2
         * @param s
         */
        bar: (s) => { },
        /**
         * comment3
         * @param s
         */
        bar2(s) { },
    };
};
exports.foo = foo;
class Foo {
    /**
     * comment4
     * @param s
     */
    bar(s) {
    }
}
exports.Foo = Foo;
/**
* comment5
*/
exports.someMethod = null.someMethod;


//// [declarationEmitRetainsJsdocyComments.d.ts]
/**
 * comment1
 * @param p
 */
export const foo: (p: string) => {
    /**
     * comment2
     * @param s
     */
    bar: (s: number) => void;
    /**
     * comment3
     * @param s
     */
    bar2(s: number): void;
};
export class Foo {
    /**
     * comment4
     * @param s
     */
    bar(s: number): void;
}
export let 
/**
* comment5
*/
someMethod: any;
global {
    interface ExtFunc {
        /**
        * comment6
        */
        someMethod(collection: any[]): boolean;
    }
}


//// [DtsFileErrors]


declarationEmitRetainsJsdocyComments.d.ts(29,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitRetainsJsdocyComments.d.ts (1 errors) ====
    /**
     * comment1
     * @param p
     */
    export const foo: (p: string) => {
        /**
         * comment2
         * @param s
         */
        bar: (s: number) => void;
        /**
         * comment3
         * @param s
         */
        bar2(s: number): void;
    };
    export class Foo {
        /**
         * comment4
         * @param s
         */
        bar(s: number): void;
    }
    export let 
    /**
    * comment5
    */
    someMethod: any;
    global {
    ~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        interface ExtFunc {
            /**
            * comment6
            */
            someMethod(collection: any[]): boolean;
        }
    }
    