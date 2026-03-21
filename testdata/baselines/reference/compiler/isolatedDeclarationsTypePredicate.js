//// [tests/cases/compiler/isolatedDeclarationsTypePredicate.ts] ////

//// [isolatedDeclarationsTypePredicate.ts]
export function isString(value: unknown) {
  return typeof value === "string";
}


//// [isolatedDeclarationsTypePredicate.js]
export function isString(value) {
    return typeof value === "string";
}


//// [isolatedDeclarationsTypePredicate.d.ts]
export declare function isString(value: unknown): value is string;
