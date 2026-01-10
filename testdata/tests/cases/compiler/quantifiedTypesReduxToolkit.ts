// @lib: esnext
// https://github.com/microsoft/TypeScript/issues/51612#issuecomment-1375431039

declare const createSlice:
  < S extends object
  , Rs extends <T extends string> { [T1 in T]: <P> Reducer<S, T1, P> }
  >
    (slice: {
      name: string,
      initialState: S,
      reducers: Rs
    }) =>
      Slice<S, Rs>

type Reducer<S, T, P> =
  | ((state: S, action: P & { type: T }) => void)
  | { reducer: (state: S, action: P & { type: T }) => void
    , prepare: (...a: never) => P
    }

type Slice<S, Rs> = 
  { actions:
      { [K in keyof Rs]:
          Rs[K] extends { prepare: (...a: infer A) => infer R } ? (...a: A) => { type: K } & R :
          Rs[K] extends (state: never, action: PayloadAction<infer P>) => unknown ? (payload: P) => { type: K, payload: P } :
          never
      }
  }

type PayloadAction<P> =
  { type: string
  , payload: P
  }

const slice = createSlice({
  name: "someSlice",
  initialState: {
    foo: "bar"
  },
  reducers: {
    simpleReducer: (state, action: PayloadAction<string>) => {
      state.foo = action.payload
    },
    reducerWithPrepareNotation: {
      prepare: (char: string, repeats: number) => {
        return { payload: char.repeat(repeats), extraStuff: true }
      },
      reducer: (state, action) => {
        state.foo = action.payload
      }
    },
    reducerWithAnotherPrepareNotation: {
      prepare: (char: string, repeats: number) => {
        return { payload: repeats * char.length }
      },
      reducer: (state, action) => {
        state.foo = state.foo.slice(0, action.payload)
      },
    },
    invalidReducerWithPrepareNotation: {
      prepare: (char: string, repeats: number) => {
        return { payload: repeats * char.length }
      },
      reducer: (state, action: PayloadAction<string>) => {
        state.foo = action.payload
      },
    },
  }
})

{
  const _expectType: (payload: string) => PayloadAction<string> = slice.actions.simpleReducer
}
{
  const _expectType: (char: string, repeats: number) => PayloadAction<string> = slice.actions.reducerWithPrepareNotation
}
{
  const _expectType: (char: string, repeats: number) => PayloadAction<number> = slice.actions.reducerWithAnotherPrepareNotation
}
