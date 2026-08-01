// @allowJs: true
// @checkJs: true
// @noEmit: true

// @filename: factory.js
/** @typedef {{ enabled: boolean }} Options */
/** @returns {void} */
module.exports = () => {};

// @filename: index.d.ts
export type Options = import("./factory").Options;
