//// [tests/cases/compiler/jsdocTypedefDeclarationComment.ts] ////

//// [main.js]
export const value = 0;

/** Use `@typedef` when documenting types. */
export const documented = 1;

/** Comment on the `Inline` type @typedef {string} Inline */

/**
 * Comment on the `Foo` type
 *
 * @typedef {Object} Foo
 * @property {boolean} bool Whether `.bool` is true or not
 */




//// [main.d.ts]
export declare const value = 0;
/** Use `@typedef` when documenting types. */
export declare const documented = 1;
/**
 * Comment on the `Inline` type
 */
export type Inline = string;
/**
 * Comment on the `Foo` type
 */
export type Foo = {
    /**
     * Whether `.bool` is true or not
     */
    bool: boolean;
};
