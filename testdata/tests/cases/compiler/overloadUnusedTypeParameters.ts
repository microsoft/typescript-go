// @strict: true
// @noUnusedLocals: true
// @noUnusedParameters: true
// @target: esnext

export function tryParseJson<T>(text: string): unknown;
export function tryParseJson<T>(text: string, predicate: (parsed: unknown) => parsed is T): T | undefined;
export function tryParseJson<T>() {
    throw new Error("noop")
}

export function tryParseJson2<T>(text: string): unknown;
export function tryParseJson2<T>(text: string, predicate: (parsed: unknown) => parsed is T): T | undefined;
export function tryParseJson2() {
    throw new Error("noop")
}

export function tryParseJson3<T>(_text: string): unknown {
    throw new Error("noop")
}