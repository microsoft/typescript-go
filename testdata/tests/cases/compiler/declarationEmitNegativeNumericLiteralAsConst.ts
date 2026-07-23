// @strict: true
// @declaration: true

export const a = -1e500 as const;
export const b = -123456789012345678901234567890 as const;
export const c = -0xff as const;
export const d = -1e3 as const;
export const e = 1e3 as const;
export const f = 0xff as const;
export const g = -0xffn as const;
export const h = -1_000 as const;
export const nested = [-1e500] as const;
export const obj = { value: -1e500 } as const;
