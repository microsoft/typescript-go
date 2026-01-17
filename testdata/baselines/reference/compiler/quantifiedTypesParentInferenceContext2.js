//// [tests/cases/compiler/quantifiedTypesParentInferenceContext2.ts] ////

//// [quantifiedTypesParentInferenceContext2.ts]
declare const f: <A>(a: A, f: [<T> ((a: A) => void)]) => void
f(0, [a => {}])

//// [quantifiedTypesParentInferenceContext2.js]
f(0, [a => { }]);
