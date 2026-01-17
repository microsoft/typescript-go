//// [tests/cases/compiler/quantifiedTypesAdvanced.ts] ////

//// [quantifiedTypesAdvanced.ts]
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


let t1 = f([
  {
    a: true,
    ab: a => +a,
    bc: b => typeof b === "number"
  },
  {
    a: "hello",
    ab: a => a + " world",
    bc: b => +b,
    extra: "foo" // TODO: an extra property should be allowed
  }
])

//// [quantifiedTypesAdvanced.js]
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
            b;
            return +b;
        }
    },
    {
        a: 42,
        ab: a => a.toString()
    }
]);
let t1 = f([
    {
        a: true,
        ab: a => +a,
        bc: b => typeof b === "number"
    },
    {
        a: "hello",
        ab: a => a + " world",
        bc: b => +b,
        extra: "foo" // TODO: an extra property should be allowed
    }
]);
