// @allowJs: true
// @checkJs: true
// @declaration: true
// @emitDeclarationOnly: true
// @filename: topComment-js.js

/** Comment on nothing */

/** Comment on noop */
(function noop() {})();

export const pi = 3;

// @filename: nonTopComment-js.js
export const e = 3;

/** Comment on nothing */

/** Comment on noop */
(function noop() {})();

export const pi = 3;

// @filename: attachedComment-js.js
/** Comment on pi */
export const pi = 3;
