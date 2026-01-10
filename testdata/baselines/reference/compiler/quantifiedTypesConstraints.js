//// [tests/cases/compiler/quantifiedTypesConstraints.ts] ////

//// [quantifiedTypesConstraints.ts]
type Input<T> = { values: T[], identifier: (value: T) => string }
declare const f: (t: <T extends object> Input<T>) => void

f({
  values: ["a", "b", "c"],
  identifier: v => v
})

//// [quantifiedTypesConstraints.js]
f({
    values: ["a", "b", "c"],
    identifier: v => v
});
