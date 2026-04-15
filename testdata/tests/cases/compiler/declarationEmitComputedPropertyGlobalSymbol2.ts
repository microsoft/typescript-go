// @declaration: true
// @isolatedDeclarations: true
// @target: es2015

export const symbolNamed = {
    [Symbol.toStringTag]: "demo",
    [Symbol.iterator]() {
        return [1, 2, 3][Symbol.iterator]();
    },
} as const;
