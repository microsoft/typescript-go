//// [tests/cases/compiler/quantifiedTypesParentInferenceContext3.ts] ////

//// [quantifiedTypesParentInferenceContext3.ts]
declare const f: <A>(x: { a: A, f: [<T> ((a: A) => void)] }) => void
f({ a: 0, f: [a => { a satisfies number; }] })


//// [quantifiedTypesParentInferenceContext3.js]
f({ a: 0, f: [a => { a; }] });
