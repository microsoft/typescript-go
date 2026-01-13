//// [tests/cases/compiler/declarationEmitOptionalParams.ts] ////

//// [declarationEmitOptionalParams.ts]
// Simple optional string
export const count = (label?: string): void => { console.log(label) }

// Optional object parameter
export const fetch = (url: string, options?: { timeout: number }): void => { console.log(url, options) }

// Multiple optional params
export const multi = (a?: string, b?: number): void => { console.log(a, b) }

// Optional with union type
export const unionOptional = (value?: string | number): void => { console.log(value) }

// Rest params after optional
export const withRest = (label?: string, ...args: Array<unknown>): void => { console.log(label, args) }


//// [declarationEmitOptionalParams.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.withRest = exports.unionOptional = exports.multi = exports.fetch = exports.count = void 0;
// Simple optional string
const count = (label) => { console.log(label); }
// Optional object parameter
;
exports.count = count;
// Optional object parameter
const fetch = (url, options) => { console.log(url, options); }
// Multiple optional params
;
exports.fetch = fetch;
// Multiple optional params
const multi = (a, b) => { console.log(a, b); }
// Optional with union type
;
exports.multi = multi;
// Optional with union type
const unionOptional = (value) => { console.log(value); }
// Rest params after optional
;
exports.unionOptional = unionOptional;
// Rest params after optional
const withRest = (label, ...args) => { console.log(label, args); };
exports.withRest = withRest;


//// [declarationEmitOptionalParams.d.ts]
export declare const count: (label?: string) => void;
export declare const fetch: (url: string, options?: {
    timeout: number;
}) => void;
export declare const multi: (a?: string, b?: number) => void;
export declare const unionOptional: (value?: string | number) => void;
export declare const withRest: (label?: string, ...args: Array<unknown>) => void;
