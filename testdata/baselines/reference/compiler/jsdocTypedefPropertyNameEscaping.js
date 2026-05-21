//// [tests/cases/compiler/jsdocTypedefPropertyNameEscaping.ts] ////

//// [index.js]
/**
 * @typedef {Object} ButtonProps
 * @property {string} label The button label
 * @property {string | null | undefined} [data-test-name] Test automation attribute
 * @property {string | null | undefined} [aria-label] Accessibility label
 * @property {string | undefined} [`back-quoted-name`] Backquoted property
 */

/**
 * @param {ButtonProps} props
 * @returns {ButtonProps}
 */
export function Button(props) {
    return props;
}

/** @typedef {string} typedef-name */

/**
 * @callback callback-name
 * @param {string} data-test-name
 * @param {string} [`back-quoted-param`]
 * @returns {void}
 */

/**
 * @template template-name
 * @returns {void}
 */
export function templated() {}




//// [index.d.ts]
/**
 * @typedef {Object} ButtonProps
 * @property {string} label The button label
 * @property {string | null | undefined} [data-test-name] Test automation attribute
 * @property {string | null | undefined} [aria-label] Accessibility label
 * @property {string | undefined} [`back-quoted-name`] Backquoted property
 */
export type ButtonProps = {
    label: string;
    data-test-name?: string | null | undefined;
    aria-label?: string | null | undefined;
    back-quoted-name?: string | undefined;
};
/**
 * @param {ButtonProps} props
 * @returns {ButtonProps}
 */
export declare function Button(props: ButtonProps): ButtonProps;
export type typedef-name = string;
export type callback-name = (data-test-name: string, back-quoted-param?: string) => void;
/** @typedef {string} typedef-name */
/**
 * @callback callback-name
 * @param {string} data-test-name
 * @param {string} [`back-quoted-param`]
 * @returns {void}
 */
/**
 * @template template-name
 * @returns {void}
 */
export declare function templated<template-name>(): void;


//// [DtsFileErrors]


index.d.ts(10,5): error TS1131: Property or signature expected.
index.d.ts(10,5): error TS2304: Cannot find name 'data'.
index.d.ts(10,10): error TS2593: Cannot find name 'test'. Do you need to install type definitions for a test runner? Try `npm i --save-dev @types/jest` or `npm i --save-dev @types/mocha` and then add 'jest' or 'mocha' to the types field in your tsconfig.
index.d.ts(10,15): error TS2363: The right-hand side of an arithmetic operation must be of type 'any', 'number', 'bigint' or an enum type.
index.d.ts(10,20): error TS1109: Expression expected.
index.d.ts(10,22): error TS2693: 'string' only refers to a type, but is being used as a value here.
index.d.ts(10,31): error TS18050: The value 'null' cannot be used here.
index.d.ts(10,38): error TS18050: The value 'undefined' cannot be used here.
index.d.ts(11,5): error TS2304: Cannot find name 'aria'.
index.d.ts(11,10): error TS2304: Cannot find name 'label'.
index.d.ts(11,16): error TS1109: Expression expected.
index.d.ts(11,18): error TS2693: 'string' only refers to a type, but is being used as a value here.
index.d.ts(11,27): error TS18050: The value 'null' cannot be used here.
index.d.ts(11,34): error TS18050: The value 'undefined' cannot be used here.
index.d.ts(12,5): error TS2304: Cannot find name 'back'.
index.d.ts(12,10): error TS2304: Cannot find name 'quoted'.
index.d.ts(12,17): error TS2363: The right-hand side of an arithmetic operation must be of type 'any', 'number', 'bigint' or an enum type.
index.d.ts(12,22): error TS1109: Expression expected.
index.d.ts(12,24): error TS2693: 'string' only refers to a type, but is being used as a value here.
index.d.ts(12,33): error TS18050: The value 'undefined' cannot be used here.
index.d.ts(13,1): error TS1128: Declaration or statement expected.
index.d.ts(19,20): error TS1005: '=' expected.
index.d.ts(19,26): error TS1005: ';' expected.
index.d.ts(19,28): error TS2693: 'string' only refers to a type, but is being used as a value here.
index.d.ts(20,21): error TS1005: '=' expected.
index.d.ts(20,27): error TS1005: ';' expected.
index.d.ts(20,30): error TS2304: Cannot find name 'data'.
index.d.ts(20,35): error TS2593: Cannot find name 'test'. Do you need to install type definitions for a test runner? Try `npm i --save-dev @types/jest` or `npm i --save-dev @types/mocha` and then add 'jest' or 'mocha' to the types field in your tsconfig.
index.d.ts(20,40): error TS2363: The right-hand side of an arithmetic operation must be of type 'any', 'number', 'bigint' or an enum type.
index.d.ts(20,44): error TS1005: ')' expected.
index.d.ts(20,46): error TS2693: 'string' only refers to a type, but is being used as a value here.
index.d.ts(20,46): error TS2695: Left side of comma operator is unused and has no side effects.
index.d.ts(20,54): error TS2304: Cannot find name 'back'.
index.d.ts(20,59): error TS2304: Cannot find name 'quoted'.
index.d.ts(20,66): error TS2304: Cannot find name 'param'.
index.d.ts(20,72): error TS1109: Expression expected.
index.d.ts(20,74): error TS2693: 'string' only refers to a type, but is being used as a value here.
index.d.ts(20,80): error TS1005: ';' expected.
index.d.ts(20,82): error TS1128: Declaration or statement expected.
index.d.ts(20,89): error TS1109: Expression expected.
index.d.ts(32,25): error TS7010: 'templated', which lacks return-type annotation, implicitly has an 'any' return type.
index.d.ts(32,43): error TS1005: ',' expected.
index.d.ts(32,50): error TS1109: Expression expected.
index.d.ts(32,51): error TS1005: ';' expected.
index.d.ts(32,57): error TS1109: Expression expected.


==== index.d.ts (45 errors) ====
    /**
     * @typedef {Object} ButtonProps
     * @property {string} label The button label
     * @property {string | null | undefined} [data-test-name] Test automation attribute
     * @property {string | null | undefined} [aria-label] Accessibility label
     * @property {string | undefined} [`back-quoted-name`] Backquoted property
     */
    export type ButtonProps = {
        label: string;
        data-test-name?: string | null | undefined;
        ~~~~
!!! error TS1131: Property or signature expected.
        ~~~~
!!! error TS2304: Cannot find name 'data'.
             ~~~~
!!! error TS2593: Cannot find name 'test'. Do you need to install type definitions for a test runner? Try `npm i --save-dev @types/jest` or `npm i --save-dev @types/mocha` and then add 'jest' or 'mocha' to the types field in your tsconfig.
                  ~~~~
!!! error TS2363: The right-hand side of an arithmetic operation must be of type 'any', 'number', 'bigint' or an enum type.
                       ~
!!! error TS1109: Expression expected.
                         ~~~~~~
!!! error TS2693: 'string' only refers to a type, but is being used as a value here.
                                  ~~~~
!!! error TS18050: The value 'null' cannot be used here.
                                         ~~~~~~~~~
!!! error TS18050: The value 'undefined' cannot be used here.
        aria-label?: string | null | undefined;
        ~~~~
!!! error TS2304: Cannot find name 'aria'.
             ~~~~~
!!! error TS2304: Cannot find name 'label'.
                   ~
!!! error TS1109: Expression expected.
                     ~~~~~~
!!! error TS2693: 'string' only refers to a type, but is being used as a value here.
                              ~~~~
!!! error TS18050: The value 'null' cannot be used here.
                                     ~~~~~~~~~
!!! error TS18050: The value 'undefined' cannot be used here.
        back-quoted-name?: string | undefined;
        ~~~~
!!! error TS2304: Cannot find name 'back'.
             ~~~~~~
!!! error TS2304: Cannot find name 'quoted'.
                    ~~~~
!!! error TS2363: The right-hand side of an arithmetic operation must be of type 'any', 'number', 'bigint' or an enum type.
                         ~
!!! error TS1109: Expression expected.
                           ~~~~~~
!!! error TS2693: 'string' only refers to a type, but is being used as a value here.
                                    ~~~~~~~~~
!!! error TS18050: The value 'undefined' cannot be used here.
    };
    ~
!!! error TS1128: Declaration or statement expected.
    /**
     * @param {ButtonProps} props
     * @returns {ButtonProps}
     */
    export declare function Button(props: ButtonProps): ButtonProps;
    export type typedef-name = string;
                       ~
!!! error TS1005: '=' expected.
                             ~
!!! error TS1005: ';' expected.
                               ~~~~~~
!!! error TS2693: 'string' only refers to a type, but is being used as a value here.
    export type callback-name = (data-test-name: string, back-quoted-param?: string) => void;
                        ~
!!! error TS1005: '=' expected.
                              ~
!!! error TS1005: ';' expected.
                                 ~~~~
!!! error TS2304: Cannot find name 'data'.
                                      ~~~~
!!! error TS2593: Cannot find name 'test'. Do you need to install type definitions for a test runner? Try `npm i --save-dev @types/jest` or `npm i --save-dev @types/mocha` and then add 'jest' or 'mocha' to the types field in your tsconfig.
                                           ~~~~
!!! error TS2363: The right-hand side of an arithmetic operation must be of type 'any', 'number', 'bigint' or an enum type.
                                               ~
!!! error TS1005: ')' expected.
                                                 ~~~~~~
!!! error TS2693: 'string' only refers to a type, but is being used as a value here.
                                                 ~~~~~~
!!! error TS2695: Left side of comma operator is unused and has no side effects.
                                                         ~~~~
!!! error TS2304: Cannot find name 'back'.
                                                              ~~~~~~
!!! error TS2304: Cannot find name 'quoted'.
                                                                     ~~~~~
!!! error TS2304: Cannot find name 'param'.
                                                                           ~
!!! error TS1109: Expression expected.
                                                                             ~~~~~~
!!! error TS2693: 'string' only refers to a type, but is being used as a value here.
                                                                                   ~
!!! error TS1005: ';' expected.
                                                                                     ~~
!!! error TS1128: Declaration or statement expected.
                                                                                            ~
!!! error TS1109: Expression expected.
    /** @typedef {string} typedef-name */
    /**
     * @callback callback-name
     * @param {string} data-test-name
     * @param {string} [`back-quoted-param`]
     * @returns {void}
     */
    /**
     * @template template-name
     * @returns {void}
     */
    export declare function templated<template-name>(): void;
                            ~~~~~~~~~
!!! error TS7010: 'templated', which lacks return-type annotation, implicitly has an 'any' return type.
                                              ~
!!! error TS1005: ',' expected.
                                                     ~
!!! error TS1109: Expression expected.
                                                      ~
!!! error TS1005: ';' expected.
                                                            ~
!!! error TS1109: Expression expected.
    