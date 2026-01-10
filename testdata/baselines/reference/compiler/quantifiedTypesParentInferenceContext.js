//// [tests/cases/compiler/quantifiedTypesParentInferenceContext.ts] ////

//// [quantifiedTypesParentInferenceContext.ts]
declare const f: <A>(a: A, f: <T> [(a: A) => void]) => void
f(0, [a => {}])


//// [quantifiedTypesParentInferenceContext.js]
f(0, [a => { }]);
