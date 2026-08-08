//// [tests/cases/compiler/enumComputedPropertyDeclarationEmit.ts] ////

//// [enumComputedPropertyDeclarationEmit.ts]
export enum StringEnum {
    A = "a",
    B = "not-an-identifier",
    Unused = "unused",
}

export const stringRecord = {
    [StringEnum.A]: StringEnum.A,
    [StringEnum.B]: StringEnum.B,
} as const;

export type StringKey = keyof typeof stringRecord;

export enum NumericEnum {
    Zero = 0,
    Negative = -1,
    Unused = 1,
}

export const numericRecord = {
    [NumericEnum.Zero]: NumericEnum.Zero,
    [NumericEnum.Negative]: NumericEnum.Negative,
} as const;

export type NumericKey = keyof typeof numericRecord;

type Assignability<T extends StringEnum | NumericEnum> = T;
export type StringDemo<T extends StringKey> = Assignability<T>;
export type NumericDemo<T extends NumericKey> = Assignability<T>;


//// [enumComputedPropertyDeclarationEmit.js]
export var StringEnum;
(function (StringEnum) {
    StringEnum["A"] = "a";
    StringEnum["B"] = "not-an-identifier";
    StringEnum["Unused"] = "unused";
})(StringEnum || (StringEnum = {}));
export const stringRecord = {
    [StringEnum.A]: StringEnum.A,
    [StringEnum.B]: StringEnum.B,
};
export var NumericEnum;
(function (NumericEnum) {
    NumericEnum[NumericEnum["Zero"] = 0] = "Zero";
    NumericEnum[NumericEnum["Negative"] = -1] = "Negative";
    NumericEnum[NumericEnum["Unused"] = 1] = "Unused";
})(NumericEnum || (NumericEnum = {}));
export const numericRecord = {
    [NumericEnum.Zero]: NumericEnum.Zero,
    [NumericEnum.Negative]: NumericEnum.Negative,
};


//// [enumComputedPropertyDeclarationEmit.d.ts]
export declare enum StringEnum {
    A = "a",
    B = "not-an-identifier",
    Unused = "unused"
}
export declare const stringRecord: {
    readonly [StringEnum.A]: StringEnum.A;
    readonly [StringEnum.B]: StringEnum.B;
};
export type StringKey = keyof typeof stringRecord;
export declare enum NumericEnum {
    Zero = 0,
    Negative = -1,
    Unused = 1
}
export declare const numericRecord: {
    readonly [NumericEnum.Zero]: NumericEnum.Zero;
    readonly [NumericEnum.Negative]: NumericEnum.Negative;
};
export type NumericKey = keyof typeof numericRecord;
type Assignability<T extends StringEnum | NumericEnum> = T;
export type StringDemo<T extends StringKey> = Assignability<T>;
export type NumericDemo<T extends NumericKey> = Assignability<T>;
export {};
