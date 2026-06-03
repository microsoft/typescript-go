//// [tests/cases/compiler/disallowUnerasableAssertion.ts] ////

//// [disallowUnerasableAssertion.ts]
export const c1 = 1 + 1 as number / 2;
export const c2 = (1 + 1 as number) / 2;
export const c3 = 1 + 1 as number === 2;


//// [disallowUnerasableAssertion.js]
export const c1 = 1 + 1;
/ 2;
export const c2 = (1 + 1) / 2;
export const c3 = 1 + 1 === 2;
