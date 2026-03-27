//// [tests/cases/compiler/declarationEmitComputedPropertyGlobalSymbol1.ts] ////

//// [declarationEmitComputedPropertyGlobalSymbol1.ts]
export const symbolNamed = {
    [Symbol.toStringTag]: "demo",
    [Symbol.iterator](): IterableIterator<number> {
        return [1, 2, 3][Symbol.iterator]();
    },
} as const;


//// [declarationEmitComputedPropertyGlobalSymbol1.js]
export const symbolNamed = {
    [Symbol.toStringTag]: "demo",
    [Symbol.iterator]() {
        return [1, 2, 3][Symbol.iterator]();
    },
};


//// [declarationEmitComputedPropertyGlobalSymbol1.d.ts]
export declare const symbolNamed: {
    readonly [Symbol.toStringTag]: "demo";
    readonly [Symbol.iterator]: () => IterableIterator<number>;
};
