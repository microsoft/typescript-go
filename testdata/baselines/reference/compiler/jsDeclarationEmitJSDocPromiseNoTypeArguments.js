//// [tests/cases/compiler/jsDeclarationEmitJSDocPromiseNoTypeArguments.ts] ////

//// [a.js]
/** @returns {Promise} */
export async function foo() {}




//// [a.d.ts]
/** @returns {Promise} */
export declare function foo(): Promise<any>;
