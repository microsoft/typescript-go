//// [tests/cases/compiler/quantifiedTypesUnwrapBaseType.ts] ////

//// [quantifiedTypesUnwrapBaseType.ts]
const f = (x: <T> [T, NoInfer<T>]) => {
  let [t0, t1] = x
}

const g = (f: <T> ((t: T) => T)) => {
  f(0)
}


//// [quantifiedTypesUnwrapBaseType.js]
const f = (x) => {
    let [t0, t1] = x;
};
const g = (f) => {
    f(0);
};
