//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsOptionalTypeLiteralProps2.ts] ////

//// [foo.js]
/**
 * foo
 *
 * @public
 * @param {object} opts
 * @param {number} opts.a
 * @param {number} [opts.b]
 * @param {number} [opts.c]
 * @returns {number}
 */
function foo({ a, b, c }) {
    return a + (b ?? 0) + (c ?? 0);
}


//// [foo.js]
"use strict";
/**
 * foo
 *
 * @public
 * @param {object} opts
 * @param {number} opts.a
 * @param {number} [opts.b]
 * @param {number} [opts.c]
 * @returns {number}
 */
function foo({ a, b, c }) {
    return a + (b ?? 0) + (c ?? 0);
}


//// [foo.d.ts]
/**
 * foo
 *
 * @public
 * @param {object} opts
 * @param {number} opts.a
 * @param {number} [opts.b]
 * @param {number} [opts.c]
 * @returns {number}
 */
function foo({ a, b, c }: {
    a: number;
    b?: number;
    c?: number;
}): number;


//// [DtsFileErrors]


out/foo.d.ts(11,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/foo.d.ts (1 errors) ====
    /**
     * foo
     *
     * @public
     * @param {object} opts
     * @param {number} opts.a
     * @param {number} [opts.b]
     * @param {number} [opts.c]
     * @returns {number}
     */
    function foo({ a, b, c }: {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        a: number;
        b?: number;
        c?: number;
    }): number;
    