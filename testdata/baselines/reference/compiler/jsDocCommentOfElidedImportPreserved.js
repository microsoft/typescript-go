//// [tests/cases/compiler/jsDocCommentOfElidedImportPreserved.ts] ////

//// [index.ts]
export interface Foo {}

//// [main.ts]
/**
 * Some random docs not related to foo
 */
/* trigger */
import * as x from './index.js';
export const foo = 1;


//// [index.js]
export {};
//// [main.js]
export const foo = 1;


//// [index.d.ts]
export interface Foo {
}
//// [main.d.ts]
export declare const foo = 1;
