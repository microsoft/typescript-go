// @declaration: true
// @strict: true

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