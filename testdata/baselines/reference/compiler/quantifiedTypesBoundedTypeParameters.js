//// [tests/cases/compiler/quantifiedTypesBoundedTypeParameters.ts] ////

//// [quantifiedTypesBoundedTypeParameters.ts]
type F = <T> { v: T, f: (v: T) => void };
declare let f1: F
declare let f2: F

f1.f(f1.v)

f1.f(f2.v)



/*
declare const f: (a: <T> {
  produce: (arg: string) => T,
  consume: (arg: T) => unknown,
  prop: <U> {
    produce: (arg: T) => U,
    consume: (arg: U) => unknown,
  }
}) => void

f({
  produce: a => Number(a),
  consume: x => {
    x satisfies number
  },
  prop: {
    produce: x => {
      x satisfies number
      return Boolean(x)
    },
    consume: y => {
      y satisfies boolean
    }
  }
})
*/

/*
declare const smallest:
  <T extends OneOf<number | string>>(xs: T[]) => T

smallest([1, 2, 3])
smallest(["a", "b", "c"])

smallest([1, "b", 3])
smallest([1, 2, undefined])

type OneOf<Spec> =
  <Actual> (
    Actual extends Spec
      ? IsUnit<Actual> extends true ? Actual : never
      : Spec
  )

type IsUnit<U> =
  [U] extends (U extends unknown ? [U] : never) ? true : false
*/

//// [quantifiedTypesBoundedTypeParameters.js]
f1.f(f1.v);
f1.f(f2.v);
/*
declare const f: (a: <T> {
  produce: (arg: string) => T,
  consume: (arg: T) => unknown,
  prop: <U> {
    produce: (arg: T) => U,
    consume: (arg: U) => unknown,
  }
}) => void

f({
  produce: a => Number(a),
  consume: x => {
    x satisfies number
  },
  prop: {
    produce: x => {
      x satisfies number
      return Boolean(x)
    },
    consume: y => {
      y satisfies boolean
    }
  }
})
*/
/*
declare const smallest:
  <T extends OneOf<number | string>>(xs: T[]) => T

smallest([1, 2, 3])
smallest(["a", "b", "c"])

smallest([1, "b", 3])
smallest([1, 2, undefined])

type OneOf<Spec> =
  <Actual> (
    Actual extends Spec
      ? IsUnit<Actual> extends true ? Actual : never
      : Spec
  )

type IsUnit<U> =
  [U] extends (U extends unknown ? [U] : never) ? true : false
*/ 
