//// [tests/cases/compiler/jsxMultiplePragmas.tsx] ////

//// [jsxMultiplePragmas.tsx]
/** @jsx h */
/** @jsx g */
declare const h: any, g: any;
export const x = <div/>;


//// [jsxMultiplePragmas.js]
export const x = h("div", null);
