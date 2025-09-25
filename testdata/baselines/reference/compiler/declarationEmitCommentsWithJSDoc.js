//// [tests/cases/compiler/declarationEmitCommentsWithJSDoc.ts] ////

//// [declarationEmitCommentsWithJSDoc.ts]
// Regular comment - should be removed
/**
 * JSDoc comment - should be preserved
 */
export class DbObject {
    // Regular comment - should be removed
    /**
     * JSDoc property comment
     */
    id: string = ""; // Trailing comment - should be removed
    
    // Regular comment - should be removed
    /**
     * JSDoc method comment
     * @returns void
     */
    method() { }
}

//// [declarationEmitCommentsWithJSDoc.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DbObject = void 0;
// Regular comment - should be removed
/**
 * JSDoc comment - should be preserved
 */
class DbObject {
    // Regular comment - should be removed
    /**
     * JSDoc property comment
     */
    id = ""; // Trailing comment - should be removed
    // Regular comment - should be removed
    /**
     * JSDoc method comment
     * @returns void
     */
    method() { }
}
exports.DbObject = DbObject;


//// [declarationEmitCommentsWithJSDoc.d.ts]
/**
 * JSDoc comment - should be preserved
 */
export declare class DbObject {
    /**
     * JSDoc property comment
     */
    id: string;
    /**
     * JSDoc method comment
     * @returns void
     */
    method(): void;
}
