//// [tests/cases/compiler/isolatedDeclarationsShadowedSymbolComputedName.ts] ////

//// [global.ts]
export const globalSymbol = {
    [Symbol.iterator]: 1,
};

//// [shadowed.ts]
const Symbol = {
    iterator: "iterator",
};

export const shadowedSymbol = {
    [Symbol.iterator]: 1,
};

//// [symbols.ts]
export const symbolValue = {
    iterator: "iterator",
};

//// [aliased.ts]
import { symbolValue as Symbol } from "./symbols.js";

export const aliasedSymbol = {
    [Symbol.iterator]: 1,
};




//// [global.d.ts]
export declare const globalSymbol: {
    [Symbol.iterator]: number;
};
//// [symbols.d.ts]
export declare const symbolValue: {
    iterator: string;
};
