//// [tests/cases/compiler/contextuallyTypedSymbolNamedProperties.ts] ////

//// [contextuallyTypedSymbolNamedProperties.ts]
// Repros from #43628

const A = Symbol("A");
const B = Symbol("B");

type Action =
    | {type: typeof A, data: string}
    | {type: typeof B, data: number}

declare const ab: Action;

declare function f<T extends { type: string | symbol }>(action: T, blah: { [K in T['type']]: (p: K) => void }): any;

f(ab, {
    [A]: ap => { ap.description },
    [B]: bp => { bp.description },
})

const x: { [sym: symbol]: (p: string) => void } = { [A]: s => s.length };


//// [contextuallyTypedSymbolNamedProperties.js]
"use strict";
// Repros from #43628
const A = Symbol("A");
const B = Symbol("B");
f(ab, {
    [A]: ap => { ap.description; },
    [B]: bp => { bp.description; },
});
const x = { [A]: s => s.length };


//// [contextuallyTypedSymbolNamedProperties.d.ts]
const A: unique symbol;
const B: unique symbol;
type Action = {
    type: typeof A;
    data: string;
} | {
    type: typeof B;
    data: number;
};
const ab: Action;
function f<T extends {
    type: string | symbol;
}>(action: T, blah: {
    [K in T['type']]: (p: K) => void;
}): any;
const x: {
    [sym: symbol]: (p: string) => void;
};


//// [DtsFileErrors]


contextuallyTypedSymbolNamedProperties.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== contextuallyTypedSymbolNamedProperties.d.ts (1 errors) ====
    const A: unique symbol;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    const B: unique symbol;
    type Action = {
        type: typeof A;
        data: string;
    } | {
        type: typeof B;
        data: number;
    };
    const ab: Action;
    function f<T extends {
        type: string | symbol;
    }>(action: T, blah: {
        [K in T['type']]: (p: K) => void;
    }): any;
    const x: {
        [sym: symbol]: (p: string) => void;
    };
    