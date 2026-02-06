// @declaration: true
// @strict: true
// @noTypesAndSymbols: true

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
