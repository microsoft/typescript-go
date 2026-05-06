//// [tests/cases/conformance/jsdoc/jsdocDestructuringParameterDeclaration.ts] ////

//// [a.js]
/**
 * @param {{ a: number; b: string }} args
 */
function f({ a, b }) {}




//// [a.d.ts]
/**
 * @param {{ a: number; b: string }} args
 */
function f({ a, b }: {
    a: number;
    b: string;
}): void;


//// [DtsFileErrors]


/a.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== /a.d.ts (1 errors) ====
    /**
     * @param {{ a: number; b: string }} args
     */
    function f({ a, b }: {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        a: number;
        b: string;
    }): void;
    