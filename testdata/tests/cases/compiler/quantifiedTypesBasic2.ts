declare const f1: (t: <T> T) => void
f1("hello")

declare const f2: (t: <T> { values: T[], identifier: (value: T) => string }) => void
f2({
  values: [{ key: "a" }, { key: "b" }, { key: "c" }],
  identifier: v => v.key
})
f2({
  values: [{ key: "a" }, { key: "b" }, { key: 0 }],
  identifier: v => v.key
})
