// @declaration: true
// @target: esnext

// https://github.com/microsoft/typescript-go/issues/4116
// Negative numeric literals whose values do not round-trip through their
// normalized text must keep their original source text in declaration emit.
export const a = -1e500 as const;
export const b = -123456789012345678901234567890 as const;
export const c = 1e500 as const;
export const d = -5 as const;
export const e = -0x10 as const;
export const f = -1_000_000_000_000_000_000_000_000 as const;
export const big = -123n as const;
export function fn() { return -1e500 as const; }
export const arrow = () => -1e500 as const;
export const withParam = (p = -1e500 as const) => 0;
// Members of object and tuple types print normalized values, matching TypeScript 6.0.
export const obj = { n: -1e500 } as const;
export const arr = [-1e500] as const;
