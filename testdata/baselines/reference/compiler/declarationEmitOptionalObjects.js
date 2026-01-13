//// [tests/cases/compiler/declarationEmitOptionalObjects.ts] ////

//// [declarationEmitOptionalObjects.ts]
interface Options {
  readonly timeout?: number
  readonly retries?: number
}

// Named interface optional param
export const configure = (options?: Options): void => { console.log(options) }

// Inline object type optional param
export const fetchData = (options?: { timeout: number }): void => { console.log(options) }

// Complex inline object with optional properties
export const createQueue = (options?: {
  readonly strategy?: "sliding" | "dropping" | "suspend"
}): void => { console.log(options) }

// Return type with optional
export const getConfig = (): Options | undefined => undefined

// Optional in generic position
export const withDefault = <T>(value: T, defaultValue?: T): T => value ?? defaultValue!


//// [declarationEmitOptionalObjects.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.withDefault = exports.getConfig = exports.createQueue = exports.fetchData = exports.configure = void 0;
// Named interface optional param
const configure = (options) => { console.log(options); }
// Inline object type optional param
;
exports.configure = configure;
// Inline object type optional param
const fetchData = (options) => { console.log(options); }
// Complex inline object with optional properties
;
exports.fetchData = fetchData;
// Complex inline object with optional properties
const createQueue = (options) => { console.log(options); }
// Return type with optional
;
exports.createQueue = createQueue;
// Return type with optional
const getConfig = () => undefined
// Optional in generic position
;
exports.getConfig = getConfig;
// Optional in generic position
const withDefault = (value, defaultValue) => value !== null && value !== void 0 ? value : defaultValue;
exports.withDefault = withDefault;


//// [declarationEmitOptionalObjects.d.ts]
interface Options {
    readonly timeout?: number;
    readonly retries?: number;
}
export declare const configure: (options?: Options) => void;
export declare const fetchData: (options?: {
    timeout: number;
}) => void;
export declare const createQueue: (options?: {
    readonly strategy?: "sliding" | "dropping" | "suspend";
}) => void;
export declare const getConfig: () => Options | undefined;
export declare const withDefault: <T>(value: T, defaultValue?: T) => T;
export {};
