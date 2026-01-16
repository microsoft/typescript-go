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