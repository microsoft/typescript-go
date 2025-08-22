//// [tests/cases/compiler/unusedTypeParameterCorrectBehavior.ts] ////

//// [unusedTypeParameterCorrectBehavior.ts]
export function wowee<T>() {
  throw new Error("TODO");
}

//// [unusedTypeParameterCorrectBehavior.js]
export function wowee() {
    throw new Error("TODO");
}
