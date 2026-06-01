//// [tests/cases/compiler/declarationEmitArrowParameterComments.ts] ////

//// [declarationEmitArrowParameterComments.ts]
export function gn(
    foo: boolean, // comment on foo
    bar: string, // comment on bar
    buzz: number,  // comment on buzz
): void {}

export const fn = (
    foo: boolean, // comment on foo
    bar: string, // comment on bar
    buzz: number,  // comment on buzz
) => {};




//// [declarationEmitArrowParameterComments.d.ts]
export declare function gn(foo: boolean, // comment on foo
bar: string, // comment on bar
buzz: number): void;
export declare const fn: (foo: boolean, // comment on foo
bar: string, // comment on bar
buzz: number) => void;
