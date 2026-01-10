//// [tests/cases/compiler/quantifiedTypesIntermediate.ts] ////

//// [quantifiedTypesIntermediate.ts]
type Input<T> = { values: T[], identifier: (value: T) => string }
declare const f: (t: (<T> Input<T>)[]) => void

f([
  {
    values: [{ key: "a" }, { key: "b" }, { key: "c" }],
    identifier: v => v.key
  }
])

f([
  {
    values: [{ key: "a" }, { key: "b" }, { key: "c" }],
    identifier: v => Number(v.key)
  }
])


f([
  {
    values: [{ key: "a" }, { key: "b" }, { key: "c" }],
    identifier: v => v.key
  },
  {
    values: [{ key: "a" }, { key: "b" }, { key: "c" }],
    identifier: v => Number(v.key)
  }
])


f([
  {
    values: [{ key: "a" }, { key: "b" }, { key: "c" }],
    identifier: v => v.key
  },
  {
    values: [{ id: 1 }, { id: 2 }, { id: 3 }],
    identifier: v => v.id.toString()
  }
])


//// [quantifiedTypesIntermediate.js]
f([
    {
        values: [{ key: "a" }, { key: "b" }, { key: "c" }],
        identifier: v => v.key
    }
]);
f([
    {
        values: [{ key: "a" }, { key: "b" }, { key: "c" }],
        identifier: v => Number(v.key)
    }
]);
f([
    {
        values: [{ key: "a" }, { key: "b" }, { key: "c" }],
        identifier: v => v.key
    },
    {
        values: [{ key: "a" }, { key: "b" }, { key: "c" }],
        identifier: v => Number(v.key)
    }
]);
f([
    {
        values: [{ key: "a" }, { key: "b" }, { key: "c" }],
        identifier: v => v.key
    },
    {
        values: [{ id: 1 }, { id: 2 }, { id: 3 }],
        identifier: v => v.id.toString()
    }
]);
