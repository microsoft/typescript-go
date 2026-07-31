// @strict: true
// @noEmit: true

// https://github.com/microsoft/typescript-go/issues/4581

// @filename: Table.tsx
import { MaterialReactTable, type MRT_TableOptions } from "./mrt";

interface TableProps<T> extends MRT_TableOptions<T> {}

export const Table = <T,>({ defaultColumn, ...props }: TableProps<T>) =>
    MaterialReactTable({ defaultColumn, ...props });

// @filename: mrt.d.ts
type Key<T, Done = false> = Done extends true ? never : unknown extends T ? string : T extends (T extends {} ? T : never) ? never : string;

export type MRT_TableOptions<T> = { defaultColumn: (value: T) => void } & {
    defaultColumn: { Cell: <Done = false>(options: MRT_TableOptions<T>) => Key<T, Done> };
};

export declare function MaterialReactTable<T>(props: MRT_TableOptions<T>): void;
