//// [tests/cases/compiler/declarationEmitNegativeNumericLiteralAsConst.ts] ////

//// [declarationEmitNegativeNumericLiteralAsConst.ts]
// https://github.com/microsoft/typescript-go/issues/4116
// Negative numeric literals whose values do not round-trip through their
// normalized text must keep their original source text in declaration emit.
export const a = -1e500 as const;
export const b = -123456789012345678901234567890 as const;
export const c = 1e500 as const;
export const d = -5 as const;
export const e = -0x10 as const;
export const f = -1_000_000_000_000_000_000_000_000 as const;
export const big = -123n as const;
export function fn() { return -1e500 as const; }
export const arrow = () => -1e500 as const;
export const withParam = (p = -1e500 as const) => 0;
// Members of object and tuple types print normalized values, matching TypeScript 6.0.
export const obj = { n: -1e500 } as const;
export const arr = [-1e500] as const;


//// [declarationEmitNegativeNumericLiteralAsConst.js]
// https://github.com/microsoft/typescript-go/issues/4116
// Negative numeric literals whose values do not round-trip through their
// normalized text must keep their original source text in declaration emit.
export const a = -1e500;
export const b = -123456789012345678901234567890;
export const c = 1e500;
export const d = -5;
export const e = -0x10;
export const f = -1e+24;
export const big = -123n;
export function fn() { return -1e500; }
export const arrow = () => -1e500;
export const withParam = (p = -1e500) => 0;
// Members of object and tuple types print normalized values, matching TypeScript 6.0.
export const obj = { n: -1e500 };
export const arr = [-1e500];


//// [declarationEmitNegativeNumericLiteralAsConst.d.ts]
export declare const a: -1e500;
export declare const b: -123456789012345678901234567890;
export declare const c: Infinity;
export declare const d: -5;
export declare const e: -0x10;
export declare const f: -1e+24;
export declare const big: -123n;
export declare function fn(): -1e500;
export declare const arrow: () => -1e500;
export declare const withParam: (p?: -1e500) => number;
export declare const obj: {
    readonly n: -Infinity;
};
export declare const arr: readonly [-Infinity];


//// [DtsFileErrors]


declarationEmitNegativeNumericLiteralAsConst.d.ts(3,25): error TS2749: 'Infinity' refers to a value, but is being used as a type here. Did you mean 'typeof Infinity'?
declarationEmitNegativeNumericLiteralAsConst.d.ts(12,17): error TS1110: Type expected.
declarationEmitNegativeNumericLiteralAsConst.d.ts(13,1): error TS1128: Declaration or statement expected.
declarationEmitNegativeNumericLiteralAsConst.d.ts(14,37): error TS1110: Type expected.
declarationEmitNegativeNumericLiteralAsConst.d.ts(14,46): error TS1005: ';' expected.


==== declarationEmitNegativeNumericLiteralAsConst.d.ts (5 errors) ====
    export declare const a: -1e500;
    export declare const b: -123456789012345678901234567890;
    export declare const c: Infinity;
                            ~~~~~~~~
!!! error TS2749: 'Infinity' refers to a value, but is being used as a type here. Did you mean 'typeof Infinity'?
    export declare const d: -5;
    export declare const e: -0x10;
    export declare const f: -1e+24;
    export declare const big: -123n;
    export declare function fn(): -1e500;
    export declare const arrow: () => -1e500;
    export declare const withParam: (p?: -1e500) => number;
    export declare const obj: {
        readonly n: -Infinity;
                    ~
!!! error TS1110: Type expected.
    };
    ~
!!! error TS1128: Declaration or statement expected.
    export declare const arr: readonly [-Infinity];
                                        ~
!!! error TS1110: Type expected.
                                                 ~
!!! error TS1005: ';' expected.
    