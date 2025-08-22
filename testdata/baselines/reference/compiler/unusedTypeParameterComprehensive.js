//// [tests/cases/compiler/unusedTypeParameterComprehensive.ts] ////

//// [unusedTypeParameterComprehensive.ts]
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

//// [unusedTypeParameterComprehensive.js]
// Unused type parameter in function (should get TS6133)
export function wowee() {
    throw new Error("TODO");
}
// Used type parameter in function (should NOT get error)
export function used(param) {
    return param;
}
// Type parameter used in constraint (should NOT get error)
export function constrained() {
    throw new Error("TODO");
}
// Multiple type parameters - one used, one unused
export function mixed(param) {
    return param;
}
