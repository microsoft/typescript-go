//// [tests/cases/compiler/spreadParameterTupleType.ts] ////

//// [spreadParameterTupleType.ts]
function f1() {
    type A = [s: string];
    type C = [...A, ...A];

    return function fn(...args: C) { } satisfies any
}

function f2() {
    type A = [a: string];
    type B = [b: string];
    type C = [c: string];
    type D = [...A, ...A, ...B, ...A, ...B, ...B, ...A, ...C];

    return function fn(...args: D) { } satisfies any;
}


//// [spreadParameterTupleType.js]
function f1() {
    return function fn(...args) { };
}
function f2() {
    return function fn(...args) { };
}


//// [spreadParameterTupleType.d.ts]
declare function f1(): (...args: [s: string, s: string]) => void;
declare function f2(): (...args: [a: string, a: string, b: string, a: string, b: string, b: string, a: string, c: string]) => void;
