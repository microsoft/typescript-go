//// [tests/cases/compiler/importValueOnlyUsedAsType.ts] ////

//// [provider.ts]
export class Value {
    data: string = "";
}
export type ValueData = { data: string };

//// [consumer.ts]
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


//// [provider.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Value = void 0;
class Value {
    data = "";
}
exports.Value = Value;
//// [consumer.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BaseProcessor = void 0;
exports.processRecord = processRecord;
function processRecord(value, // parameter type
callback) {
    // return type
    callback(value);
    return value;
}
class BaseProcessor {
    current = null;
}
exports.BaseProcessor = BaseProcessor;
