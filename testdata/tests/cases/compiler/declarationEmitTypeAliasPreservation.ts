// @declaration: true
// @strict: true
// @noTypesAndSymbols: true

export type NonEmptyArray<A> = [A, ...Array<A>]
export type NonEmptyReadonlyArray<A> = readonly [A, ...ReadonlyArray<A>]

// All of these should preserve the type alias in declaration emit
export const make = <A>(...elements: NonEmptyArray<A>): NonEmptyArray<A> => elements
export const fromArray = <A>(arr: Array<A>): Array<A> => arr
export const fromReadonly = <A>(arr: ReadonlyArray<A>): ReadonlyArray<A> => arr
export const first = <A>(arr: NonEmptyReadonlyArray<A>): A => arr[0]
export const allocate = <A>(n: number): Array<A | undefined> => new Array(n)
