//// [tests/cases/compiler/quantifiedTypesBoundedTypeParameters.ts] ////

//// [quantifiedTypesBoundedTypeParameters.ts]
type F = <T> { v: T, f: (v: T) => void };
declare let f1: F
declare let f2: F

f1.f(f1.v)
f1.f(f2.v)


//// [quantifiedTypesBoundedTypeParameters.js]
f1.f(f1.v);
f1.f(f2.v);
