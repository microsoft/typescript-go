//// [tests/cases/compiler/declarationEmitTypeofIndexedAccessNoParens.ts] ////

//// [declarationEmitTypeofIndexedAccessNoParens.ts]
export const C = { A: 1 };
export type C = typeof C[keyof typeof C];

// Parenthesized form should also round-trip
export type C2 = (typeof C)[keyof typeof C];

// Array type of a parsed typeof should also preserve source
export const arr = [C];
export type ArrAlias = typeof arr[number];




//// [declarationEmitTypeofIndexedAccessNoParens.d.ts]
export declare const C: {
    A: number;
};
export type C = typeof C[keyof typeof C];
export type C2 = (typeof C)[keyof typeof C];
export declare const arr: {
    A: number;
}[];
export type ArrAlias = typeof arr[number];
