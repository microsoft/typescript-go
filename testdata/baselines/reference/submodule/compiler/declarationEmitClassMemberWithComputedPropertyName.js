//// [tests/cases/compiler/declarationEmitClassMemberWithComputedPropertyName.ts] ////

//// [declarationEmitClassMemberWithComputedPropertyName.ts]
const k1 = Symbol();
const k2 = 'foo' as const;

const k3 = Symbol();
const k4 = 'prop' as const;

class Foo {
    static [k1](): number {
        return 1;
    }
    [k1](): string {
        return "";
    }

    static [k2]() {
        return 1;
    }
    [k2]() {
        return "";
    }

    static m1() {}
    m1() {}

    static [k3] = 1;
    [k3] = 1;

    static [k4] = 1;
    [k4] = 2;

    static p1 = 3;
    p1 = 4;
}

export const t1 = Foo[k1];
export const t2 = new Foo()[k1];

export const t3 = Foo[k2];
export const t4 = new Foo()[k2];

export const t5 = Foo.m1;
export const t6 = new Foo().m1;

export const t7 = Foo[k3];
export const t8 = new Foo()[k3];

export const t9 = Foo[k4];
export const t10 = new Foo()[k4];

export const t11 = Foo.p1;
export const t12 = new Foo().p1;




//// [declarationEmitClassMemberWithComputedPropertyName.d.ts]
declare class Foo {
    static m1(): void;
    m1(): void;
    static p1: number;
    p1: number;
}
export declare const t1: (typeof Foo)[typeof k1];
export declare const t2: () => string;
export declare const t3: typeof Foo.foo;
export declare const t4: () => string;
export declare const t5: typeof Foo.m1;
export declare const t6: () => void;
export declare const t7: number;
export declare const t8: number;
export declare const t9: number;
export declare const t10: number;
export declare const t11: number;
export declare const t12: number;
export {};


//// [DtsFileErrors]


declarationEmitClassMemberWithComputedPropertyName.d.ts(7,39): error TS2538: Type 'any' cannot be used as an index type.
declarationEmitClassMemberWithComputedPropertyName.d.ts(7,46): error TS2304: Cannot find name 'k1'.
declarationEmitClassMemberWithComputedPropertyName.d.ts(9,37): error TS2339: Property 'foo' does not exist on type 'typeof Foo'.


==== declarationEmitClassMemberWithComputedPropertyName.d.ts (3 errors) ====
    declare class Foo {
        static m1(): void;
        m1(): void;
        static p1: number;
        p1: number;
    }
    export declare const t1: (typeof Foo)[typeof k1];
                                          ~~~~~~~~~
!!! error TS2538: Type 'any' cannot be used as an index type.
                                                 ~~
!!! error TS2304: Cannot find name 'k1'.
    export declare const t2: () => string;
    export declare const t3: typeof Foo.foo;
                                        ~~~
!!! error TS2339: Property 'foo' does not exist on type 'typeof Foo'.
    export declare const t4: () => string;
    export declare const t5: typeof Foo.m1;
    export declare const t6: () => void;
    export declare const t7: number;
    export declare const t8: number;
    export declare const t9: number;
    export declare const t10: number;
    export declare const t11: number;
    export declare const t12: number;
    export {};
    