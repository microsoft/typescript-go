// @strict: true
// @noUnusedLocals: true
// @noUnusedParameters: true
// @target: esnext

// Unused type parameter in function (should get TS6133)
export function wowee<T>() {
  throw new Error("TODO");
}

// Used type parameter in function (should NOT get error)
export function used<T>(param: T): T {
  return param;
}

// Type parameter used in constraint (should NOT get error)
export function constrained<T extends string>(): T {
  throw new Error("TODO");
}

// Multiple type parameters - one used, one unused
export function mixed<T, U>(param: T): T {
  return param;
}