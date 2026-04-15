// @declaration: true
// @isolatedDeclarations: true
// @target: es2015

export const symbolNamed = {
    [Symbol.toStringTag]: "demo",
    [Symbol.iterator](): IterableIterator<number> {
        return [1, 2, 3][Symbol.iterator]();
    },
} as const;
