// @declaration: true
// @strict: true
// @noTypesAndSymbols: true

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
