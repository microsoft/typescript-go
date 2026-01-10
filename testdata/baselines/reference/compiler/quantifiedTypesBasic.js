//// [tests/cases/compiler/quantifiedTypesBasic.ts] ////

//// [quantifiedTypesBasic.ts]
let t0: <T> T = "hello"

let t1: <T> { values: T[], identifier: (value: T) => string } = {
  values: [{ key: "a" }, { key: "b" }, { key: "c" }],
  identifier: v => v.key
}

let t2: <T> { values: T[], identifier: (value: T) => string } = {
  values: [{ key: "a" }, { key: "b" }, { key: 0 }],
  identifier: v => v.key
}


//// [quantifiedTypesBasic.js]
let t0 = "hello";
let t1 = {
    values: [{ key: "a" }, { key: "b" }, { key: "c" }],
    identifier: v => v.key
};
let t2 = {
    values: [{ key: "a" }, { key: "b" }, { key: 0 }],
    identifier: v => v.key
};
