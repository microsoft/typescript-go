//// [tests/cases/compiler/quantifiedTypesParentInferenceContext4.ts] ////

//// [quantifiedTypesParentInferenceContext4.ts]
declare const f:
  (x:
    <A, B, D, E> {
      a: () => A,
      ab: (a: A) => B,
      nested: <C> {
        c: (b: B) => C,
        cd: (c: C) => D,
      },
      de: (d: D) => E 
    }
  ) => void

// TODO: this should compile just like the f1 call
f({
  a: () => ({ a: 0 }),
  ab: x => ({ ...x, b: "" }),
  nested: {
    c: x => ({ ...x, c: +x.b }),
    cd: x => ({ ...x, d: Boolean(x.c) })
  },
  de: x => ({ ...x, e: "" })
})

declare const f1:
  <A, B, C, D, E>
  (x:
     {
      a: () => A,
      ab: (a: A) => B,
      nested: {
        c: (b: B) => C,
        cd: (c: C) => D,
      },
      de: (d: D) => E 
    }
  ) => void

f1({
  a: () => ({ a: 0 }),
  ab: x => ({ ...x, b: "" }),
  nested: {
    c: x => ({ ...x, c: +x.b }),
    cd: x => ({ ...x, d: Boolean(x.c) })
  },
  de: x => ({ ...x, e: "" })
})

//// [quantifiedTypesParentInferenceContext4.js]
// TODO: this should compile just like the f1 call
f({
    a: () => ({ a: 0 }),
    ab: x => (Object.assign(Object.assign({}, x), { b: "" })),
    nested: {
        c: x => (Object.assign(Object.assign({}, x), { c: +x.b })),
        cd: x => (Object.assign(Object.assign({}, x), { d: Boolean(x.c) }))
    },
    de: x => (Object.assign(Object.assign({}, x), { e: "" }))
});
f1({
    a: () => ({ a: 0 }),
    ab: x => (Object.assign(Object.assign({}, x), { b: "" })),
    nested: {
        c: x => (Object.assign(Object.assign({}, x), { c: +x.b })),
        cd: x => (Object.assign(Object.assign({}, x), { d: Boolean(x.c) }))
    },
    de: x => (Object.assign(Object.assign({}, x), { e: "" }))
});
