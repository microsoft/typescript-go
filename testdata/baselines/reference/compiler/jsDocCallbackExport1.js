//// [tests/cases/compiler/jsDocCallbackExport1.ts] ////

//// [x.js]
/**
 * @callback Foo
 * @param {string} x
 * @returns {number}
 */
function f1() {}




//// [x.d.ts]
type Foo = (x: string) => number;
/**
 * @callback Foo
 * @param {string} x
 * @returns {number}
 */
function f1(): void;


//// [DtsFileErrors]


x.d.ts(7,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== x.d.ts (1 errors) ====
    type Foo = (x: string) => number;
    /**
     * @callback Foo
     * @param {string} x
     * @returns {number}
     */
    function f1(): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    