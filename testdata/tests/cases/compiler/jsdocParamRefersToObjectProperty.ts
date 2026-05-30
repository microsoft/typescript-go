// @allowJs: true
// @checkJs: true
// @noEmit: true

// @filename: main.js
export const map = /** @type {const} */ ({
    foo: 'foo',
    bar: 'bar',
});

/**
 * @param {map.foo} foo
 */
export function buzz(foo) {}
