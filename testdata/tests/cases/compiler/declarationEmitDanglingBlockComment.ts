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

// @filename: elidedDeclarationComment-js.js
/** Comment on unused */
const unused = 1;

export const pi = 3;

// @filename: topComment-exports.js
/** Comment on nothing */

/** Comment on noop */
(function noop() {})();

exports.pi = 3;

// @filename: topComment-moduleExports.js
/** Comment on nothing */

/** Comment on noop */
(function noop() {})();

module.exports = 3;

// @filename: topComment-ts.ts

/** Comment on nothing */

/** Comment on noop */
(function noop() {})();

export const pi = 3;

// @filename: nonTopComment-ts.ts
export const e = 3;

/** Comment on nothing */

/** Comment on noop */
(function noop() {})();

export const pi = 3;

// @filename: attachedComment-ts.ts
/** Comment on pi */
export const pi = 3;

// @filename: elidedDeclarationComment-ts.ts
/** Comment on unused */
const unused = 1;

export const pi = 3;
