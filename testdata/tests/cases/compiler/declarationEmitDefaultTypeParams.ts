// @declaration: true
// @noTypesAndSymbols: true

export interface Effect<A, E = never, R = never> {
  readonly _A: A
  readonly _E: E
  readonly _R: R
}

export interface Deferred<A, E = never> {
  readonly _A: A
  readonly _E: E
}

// When return type uses defaults, they should be omitted
export const succeed = <A>(value: A): Effect<A> => ({ _A: value, _E: undefined as never, _R: undefined as never })
export const fail = <E>(error: E): Effect<never, E> => ({ _A: undefined as never, _E: error, _R: undefined as never })
export const makeDeferred = <A>(): Deferred<A> => ({ _A: undefined as never, _E: undefined as never })

// Many default params - only non-defaults should appear
export interface Channel<Out, Err = never, Done = void, In = unknown, InErr = unknown, InDone = unknown, Env = never> {
  readonly _Out: Out
}

export const fromString = (s: string): Channel<string> => ({ _Out: s })
