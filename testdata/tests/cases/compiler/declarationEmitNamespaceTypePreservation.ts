// @declaration: true
// @noTypesAndSymbols: true

export namespace Pull {
  export interface Halt<E> {
    readonly _tag: "Halt"
    readonly error: E
  }

  export type ExcludeHalt<E> = Exclude<E, Halt<unknown>>
}

export namespace Arr {
  export type NonEmptyReadonlyArray<A> = readonly [A, ...ReadonlyArray<A>]
}

// Namespace-qualified types should be preserved
export const excludeHalt = <E>(error: E): Pull.ExcludeHalt<E> => error as Pull.ExcludeHalt<E>
export const toNonEmpty = <A>(arr: Arr.NonEmptyReadonlyArray<A>): Arr.NonEmptyReadonlyArray<A> => arr
export const process = <E>(input: E): Pull.ExcludeHalt<E> | null => input as Pull.ExcludeHalt<E>
