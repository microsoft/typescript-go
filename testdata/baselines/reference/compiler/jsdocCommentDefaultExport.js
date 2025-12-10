//// [tests/cases/compiler/jsdocCommentDefaultExport.ts] ////

//// [jsdocCommentDefaultExport.ts]
/** Some comment */
export default {
    fn() {}
}


//// [jsdocCommentDefaultExport.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
/** Some comment */
exports.default = {
    fn() { }
};


//// [jsdocCommentDefaultExport.d.ts]
/** Some comment */
declare const _default: {
    fn(): void;
};
export default _default;
