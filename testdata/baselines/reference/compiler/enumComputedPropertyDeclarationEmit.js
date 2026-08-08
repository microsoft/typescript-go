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
    readonly a: StringEnum.A;
    readonly "not-an-identifier": StringEnum.B;
};
export type StringKey = keyof typeof stringRecord;
export declare enum NumericEnum {
    Zero = 0,
    Negative = -1,
    Unused = 1
}
export declare const numericRecord: {
    readonly 0: NumericEnum.Zero;
    readonly [-1]: NumericEnum.Negative;
};
export type NumericKey = keyof typeof numericRecord;
type Assignability<T extends StringEnum | NumericEnum> = T;
export type StringDemo<T extends StringKey> = Assignability<T>;
export type NumericDemo<T extends NumericKey> = Assignability<T>;
export {};


//// [DtsFileErrors]


enumComputedPropertyDeclarationEmit.d.ts(22,61): error TS2344: Type 'T' does not satisfy the constraint 'NumericEnum | StringEnum'.
  Type '"a" | "not-an-identifier"' is not assignable to type 'NumericEnum | StringEnum'.
    Type '"a"' is not assignable to type 'NumericEnum | StringEnum'.


==== enumComputedPropertyDeclarationEmit.d.ts (1 errors) ====
    export declare enum StringEnum {
        A = "a",
        B = "not-an-identifier",
        Unused = "unused"
    }
    export declare const stringRecord: {
        readonly a: StringEnum.A;
        readonly "not-an-identifier": StringEnum.B;
    };
    export type StringKey = keyof typeof stringRecord;
    export declare enum NumericEnum {
        Zero = 0,
        Negative = -1,
        Unused = 1
    }
    export declare const numericRecord: {
        readonly 0: NumericEnum.Zero;
        readonly [-1]: NumericEnum.Negative;
    };
    export type NumericKey = keyof typeof numericRecord;
    type Assignability<T extends StringEnum | NumericEnum> = T;
    export type StringDemo<T extends StringKey> = Assignability<T>;
                                                                ~
!!! error TS2344: Type 'T' does not satisfy the constraint 'NumericEnum | StringEnum'.
!!! error TS2344:   Type '"a" | "not-an-identifier"' is not assignable to type 'NumericEnum | StringEnum'.
!!! error TS2344:     Type '"a"' is not assignable to type 'NumericEnum | StringEnum'.
    export type NumericDemo<T extends NumericKey> = Assignability<T>;
    export {};
    