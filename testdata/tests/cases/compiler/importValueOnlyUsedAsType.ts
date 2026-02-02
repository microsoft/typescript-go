// @module: commonjs
// @target: es2020

// @filename: provider.ts
export class Value {
    data: string = "";
}
export type ValueData = { data: string };

// @filename: consumer.ts
// Import Value as a value (no `type` keyword), but only use it in type positions
import { Value, type ValueData } from "./provider";

// Value is ONLY used in type positions:
export interface Record {
    getValue(): Value; // return type
    setValue(value: Value): void; // parameter type
    readonly currentValue: Value; // type annotation
}

export function processRecord(
    value: Value, // parameter type
    callback: (result: Value) => void, // parameter type in callback
): Value {
    // return type
    callback(value);
    return value;
}

export class BaseProcessor {
    current: Value | null = null;
}
