// @strict: true
// @noEmit: true
// @skipLibCheck: true
// @module: esnext
// @moduleResolution: bundler

// @filename: lib.d.ts
export interface TableCommonProps {
    style?: string;
    className?: string;
    role?: string;
}

export interface TableKeyedProps extends TableCommonProps {
    key: string;
}

export interface TableRowProps extends TableKeyedProps {}

// @filename: main.ts
import type { TableRowProps } from "./lib";

type AriaRole =
    | "alert"
    | "button"
    | "checkbox"
    | "dialog"
    | "grid"
    | "row"
    | (string & {});

interface HTMLAttributes {
    role?: AriaRole | undefined;
    onMouseDown?: (() => void) | undefined;
}

declare module "./lib" {
    export interface TableRowProps extends Omit<HTMLAttributes, "onMouseDown"> {
        "data-test-id"?: string;
    }
}

declare const row: TableRowProps;
row.role;
