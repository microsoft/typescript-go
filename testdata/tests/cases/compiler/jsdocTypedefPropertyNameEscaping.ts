// @allowJs: true
// @checkJs: true
// @declaration: true
// @emitDeclarationOnly: true

// @filename: index.js
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
