//// [tests/cases/compiler/declarationEmitComputedPropertyGlobalSymbol2.ts] ////

//// [declarationEmitComputedPropertyGlobalSymbol2.ts]
export const symbolNamed = {
    [Symbol.toStringTag]: "demo",
    [Symbol.iterator]() {
        return [1, 2, 3][Symbol.iterator]();
    },
} as const;


//// [declarationEmitComputedPropertyGlobalSymbol2.js]
export const symbolNamed = {
    [Symbol.toStringTag]: "demo",
    [Symbol.iterator]() {
        return [1, 2, 3][Symbol.iterator]();
    },
};


//// [declarationEmitComputedPropertyGlobalSymbol2.d.ts]
export declare const symbolNamed: {
    readonly [Symbol.toStringTag]: "demo";
    readonly [Symbol.iterator]: () => ArrayIterator<number>;
};
