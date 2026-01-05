declare const f:
  <T extends (<A extends string | number, B, C> { a: A, ab: (a: A) => B, bc?: (b: B) => C })[]>(a: [...T]) =>
      { [K in keyof T]: T[K] extends { bc: (...a: never) => infer C } ? C : "lol" }

let t0 = f([
  {
    a: "0",
    ab: a => +a,
    bc: b => typeof b === "number"
  },
  {
    a: "hello",
    ab: a => a + " world",
    bc: b => {
      b satisfies string
      return +b
    }
  },
  {
    a: 42,
    ab: a => a.toString()
  }
])

