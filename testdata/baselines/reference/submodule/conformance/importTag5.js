//// [tests/cases/conformance/jsdoc/importTag5.ts] ////

//// [types.ts]
export interface Foo {
    a: number;
}

//// [foo.js]
/**
 * @import { Foo } from "./types"
 */

/**
 * @param { Foo } foo
 */
function f(foo) {}




//// [types.d.ts]
export interface Foo {
    a: number;
}
//// [foo.d.ts]
/**
 * @import { Foo } from "./types"
 */
import type { Foo } from "./types";
/**
 * @param { Foo } foo
 */
function f(foo: Foo): void;


//// [DtsFileErrors]


/foo.d.ts(8,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== /types.d.ts (0 errors) ====
    export interface Foo {
        a: number;
    }
    
==== /foo.d.ts (1 errors) ====
    /**
     * @import { Foo } from "./types"
     */
    import type { Foo } from "./types";
    /**
     * @param { Foo } foo
     */
    function f(foo: Foo): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    