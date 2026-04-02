//// [tests/cases/compiler/declarationEmitReadonlyTupleFromConstSatisfies.ts] ////

//// [declarationEmitReadonlyTupleFromConstSatisfies.ts]
// Regression test for #3192
// `as const satisfies` with a mutable array contextual type should not
// emit a readonly tuple in declaration output.

export const obj = { 
  array: [
    { n: 1 }, 
    { n: 2 },
  ],
} as const satisfies { array?: Readonly<{ n: unknown; }>[] }


//// [declarationEmitReadonlyTupleFromConstSatisfies.js]
// Regression test for #3192
// `as const satisfies` with a mutable array contextual type should not
// emit a readonly tuple in declaration output.
export const obj = {
    array: [
        { n: 1 },
        { n: 2 },
    ],
};


//// [declarationEmitReadonlyTupleFromConstSatisfies.d.ts]
export declare const obj: {
    readonly array: [{
        readonly n: 1;
    }, {
        readonly n: 2;
    }];
};
