// @allowJs: true
// @checkJs: true
// @noImplicitAny: true
// @noEmit: true

// @filename: a.js
/**
 * @overload
 * @param {number} x
 */

/**
 * @overload
 * @param {string} x
 */

/**
 * @param {string | number} x
 * @returns {string | number}
 */
function id(x) {
    return x;
}
