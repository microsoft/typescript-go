//// [tests/cases/compiler/declarationEmitDanglingBlockComment.ts] ////

//// [topComment-js.js]
/** Comment on nothing */

/** Comment on noop */
(function noop() {})();

export const pi = 3;

//// [nonTopComment-js.js]
export const e = 3;

/** Comment on nothing */

/** Comment on noop */
(function noop() {})();

export const pi = 3;

//// [attachedComment-js.js]
/** Comment on pi */
export const pi = 3;

//// [elidedDeclarationComment-js.js]
/** Comment on unused */
const unused = 1;

export const pi = 3;

//// [topComment-exports.js]
/** Comment on nothing */

/** Comment on noop */
(function noop() {})();

exports.pi = 3;

//// [topComment-moduleExports.js]
/** Comment on nothing */

/** Comment on noop */
(function noop() {})();

module.exports = 3;

//// [topComment-ts.ts]
/** Comment on nothing */

/** Comment on noop */
(function noop() {})();

export const pi = 3;

//// [nonTopComment-ts.ts]
export const e = 3;

/** Comment on nothing */

/** Comment on noop */
(function noop() {})();

export const pi = 3;

//// [attachedComment-ts.ts]
/** Comment on pi */
export const pi = 3;

//// [elidedDeclarationComment-ts.ts]
/** Comment on unused */
const unused = 1;

export const pi = 3;




//// [topComment-js.d.ts]
export declare const pi = 3;
//// [nonTopComment-js.d.ts]
export declare const e = 3;
export declare const pi = 3;
//// [attachedComment-js.d.ts]
/** Comment on pi */
export declare const pi = 3;
//// [elidedDeclarationComment-js.d.ts]
export declare const pi = 3;
//// [topComment-exports.d.ts]
export declare var pi: 3;
//// [topComment-moduleExports.d.ts]
declare const _exports = 3;
export = _exports;
//// [topComment-ts.d.ts]
/** Comment on nothing */
export declare const pi = 3;
//// [nonTopComment-ts.d.ts]
export declare const e = 3;
export declare const pi = 3;
//// [attachedComment-ts.d.ts]
/** Comment on pi */
export declare const pi = 3;
//// [elidedDeclarationComment-ts.d.ts]
export declare const pi = 3;
