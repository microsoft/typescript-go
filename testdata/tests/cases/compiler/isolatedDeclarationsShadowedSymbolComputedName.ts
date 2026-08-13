// @declaration: true
// @emitDeclarationOnly: true
// @isolatedDeclarations: true
// @target: esnext

// @filename: global.ts
export const globalSymbol = {
    [Symbol.iterator]: 1,
};

// @filename: shadowed.ts
const Symbol = {
    iterator: "iterator",
};

export const shadowedSymbol = {
    [Symbol.iterator]: 1,
};

// @filename: symbols.ts
export const symbolValue = {
    iterator: "iterator",
};

// @filename: aliased.ts
import { symbolValue as Symbol } from "./symbols.js";

export const aliasedSymbol = {
    [Symbol.iterator]: 1,
};
