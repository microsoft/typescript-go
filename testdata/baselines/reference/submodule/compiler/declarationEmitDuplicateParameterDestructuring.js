//// [tests/cases/compiler/declarationEmitDuplicateParameterDestructuring.ts] ////

//// [declarationEmitDuplicateParameterDestructuring.ts]
export const fn1 = ({ prop: a, prop: b }: { prop: number }) => a + b;

export const fn2 = ({ prop: a }: { prop: number }, { prop: b }: { prop: number }) => a + b;




//// [declarationEmitDuplicateParameterDestructuring.d.ts]
export const fn1: ({ prop: a, prop: b }: {
    prop: number;
}) => number;
export const fn2: ({ prop: a }: {
    prop: number;
}, { prop: b }: {
    prop: number;
}) => number;
