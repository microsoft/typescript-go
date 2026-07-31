//// [tests/cases/compiler/genericTypeForwardedConditional.ts] ////

//// [genericTypeForwardedConditional.ts]
// repro from #4581

type IsTuple<TMaybeTuple> = TMaybeTuple extends any[] ? TMaybeTuple : never;

type DeepKeys<TObj, TDone = false> = TDone extends true
    ? never
    : TObj extends any[] & IsTuple<TObj>
      ? never
      : TObj extends object
        ? keyof TObj | DeepKeysPrefix<TObj, keyof TObj>
        : never;

type DeepKeysPrefix<TParent, TPrefix> = TPrefix extends keyof TParent
    ? DeepKeys<TParent[TPrefix], true>
    : never;

declare function DataTable<TArg>(props: { accessorKey: DeepKeys<TArg> }): void;

export const Table = <TRow extends object>(accessorKey: DeepKeys<TRow>) =>
    DataTable({ accessorKey });


//// [genericTypeForwardedConditional.js]
// repro from #4581
export const Table = (accessorKey) => DataTable({ accessorKey });
