//// [tests/cases/compiler/declarationEmitNegativeNumericLiteralAsConst.ts] ////

//// [declarationEmitNegativeNumericLiteralAsConst.ts]
export const a = -1e500 as const;
export const b = -123456789012345678901234567890 as const;
export const c = -0xff as const;
export const d = -1e3 as const;
export const e = 1e3 as const;
export const f = 0xff as const;
export const g = -0xffn as const;
export const h = -1_000 as const;
export const nested = [-1e500] as const;
export const obj = { value: -1e500 } as const;


//// [declarationEmitNegativeNumericLiteralAsConst.js]
export const a = -1e500;
export const b = -123456789012345678901234567890;
export const c = -0xff;
export const d = -1e3;
export const e = 1e3;
export const f = 0xff;
export const g = -0xffn;
export const h = -1000;
export const nested = [-1e500];
export const obj = { value: -1e500 };


//// [declarationEmitNegativeNumericLiteralAsConst.d.ts]
export declare const a: -1e500;
export declare const b: -123456789012345678901234567890;
export declare const c: -0xff;
export declare const d: -1e3;
export declare const e: 1000;
export declare const f: 255;
export declare const g: -0xffn;
export declare const h: -1000;
export declare const nested: readonly [-Infinity];
export declare const obj: {
    readonly value: -Infinity;
};


//// [DtsFileErrors]


declarationEmitNegativeNumericLiteralAsConst.d.ts(9,40): error TS1110: Type expected.
declarationEmitNegativeNumericLiteralAsConst.d.ts(9,49): error TS1005: ';' expected.
declarationEmitNegativeNumericLiteralAsConst.d.ts(11,21): error TS1110: Type expected.
declarationEmitNegativeNumericLiteralAsConst.d.ts(12,1): error TS1128: Declaration or statement expected.


==== declarationEmitNegativeNumericLiteralAsConst.d.ts (4 errors) ====
    export declare const a: -1e500;
    export declare const b: -123456789012345678901234567890;
    export declare const c: -0xff;
    export declare const d: -1e3;
    export declare const e: 1000;
    export declare const f: 255;
    export declare const g: -0xffn;
    export declare const h: -1000;
    export declare const nested: readonly [-Infinity];
                                           ~
!!! error TS1110: Type expected.
                                                    ~
!!! error TS1005: ';' expected.
    export declare const obj: {
        readonly value: -Infinity;
                        ~
!!! error TS1110: Type expected.
    };
    ~
!!! error TS1128: Declaration or statement expected.
    