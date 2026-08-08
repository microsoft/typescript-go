//// [tests/cases/compiler/declarationEmitDanglingBlockCommentJs.ts] ////

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




//// [topComment-js.d.ts]
export declare const pi = 3;
//// [nonTopComment-js.d.ts]
export declare const e = 3;
export declare const pi = 3;
//// [attachedComment-js.d.ts]
/** Comment on pi */
export declare const pi = 3;
