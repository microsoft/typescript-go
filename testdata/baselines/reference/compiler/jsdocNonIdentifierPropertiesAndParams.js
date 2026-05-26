//// [tests/cases/compiler/jsdocNonIdentifierPropertiesAndParams.ts] ////

//// [index.js]
/**
 * @typedef {Object} ButtonProps
 * @property {string} label The button label
 * @property {string | null | undefined} [data-test-name] Test automation attribute
 * @property {string | null | undefined} [aria-label] Accessibility label
 */

/**
 * @param {ButtonProps} props
 * @returns {ButtonProps}
 */
export function Button(props) {
    return { ...props }
}

/**
 * @callback ButtonPropsCallback
 * @param {ButtonProps} [props-like]
 * @returns {ButtonProps}
 */




//// [index.d.ts]
/**
 * @typedef {Object} ButtonProps
 * @property {string} label The button label
 * @property {string | null | undefined} [data-test-name] Test automation attribute
 * @property {string | null | undefined} [aria-label] Accessibility label
 */
export type ButtonProps = {
    label: string;
    "data-test-name"?: string | null | undefined;
    "aria-label"?: string | null | undefined;
};
/**
 * @param {ButtonProps} props
 * @returns {ButtonProps}
 */
export declare function Button(props: ButtonProps): ButtonProps;
export type ButtonPropsCallback = (props_like?: ButtonProps) => ButtonProps;
/**
 * @callback ButtonPropsCallback
 * @param {ButtonProps} [props-like]
 * @returns {ButtonProps}
 */
