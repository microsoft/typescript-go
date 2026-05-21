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
    "data-test-name"?: string | null | undefined;
    "aria-label"?: string | null | undefined;
    "back-quoted-name"?: string | undefined;
};
/**
 * @param {ButtonProps} props
 * @returns {ButtonProps}
 */
export declare function Button(props: ButtonProps): ButtonProps;
export type _typedef_name = string;
export type _callback_name = (_data_test_name: string, _back_quoted_param?: string) => void;
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
export declare function templated<_template_name>(): void;
