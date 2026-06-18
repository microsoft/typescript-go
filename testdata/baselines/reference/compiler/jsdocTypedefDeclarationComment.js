//// [tests/cases/compiler/jsdocTypedefDeclarationComment.ts] ////

//// [main.js]
export const value = 0;

/**
 * Comment on the `Foo` type
 *
 * @typedef {Object} Foo
 * @property {boolean} bool Whether `.bool` is true or not
 */




//// [main.d.ts]
export declare const value = 0;
/**
 * Comment on the `Foo` type
 */
export type Foo = {
    /**
     * Whether `.bool` is true or not
     */
    bool: boolean;
};
