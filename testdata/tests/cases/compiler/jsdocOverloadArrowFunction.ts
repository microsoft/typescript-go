// @allowJs: true
// @checkJs: true
// @noEmit: true

// @filename: index.js
/** @typedef {{ text: string }} Result */

/**
 * @overload
 * @param {string} value
 * @param {false} withObject
 * @returns {string}
 */
/**
 * @overload
 * @param {string} value
 * @param {true} withObject
 * @returns {Result}
 */
/**
 * @param {string} value
 * @param {boolean} withObject
 * @returns {string | Result}
 */
const decode = (value, withObject) => withObject ? { text: value } : value;

decode("value", false).toUpperCase();
decode("value", true).text;
