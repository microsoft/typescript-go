//// [tests/cases/conformance/jsdoc/jsdocAccessibilityTagsDeclarations.ts] ////

//// [jsdocAccessibilityTagDeclarations.js]
class Protected {
    /** @protected */
    constructor(c) {
        /** @protected */
        this.c = c
    }
    /** @protected */
    m() {
        return this.c
    }
    /** @protected */
    get p() { return this.c }
    /** @protected */
    set p(value) { this.c = value }
}

class Private {
    /** @private */
    constructor(c) {
        /** @private */
        this.c = c
    }
    /** @private */
    m() {
        return this.c
    }
    /** @private */
    get p() { return this.c }
    /** @private */
    set p(value) { this.c = value }
}

// https://github.com/microsoft/TypeScript/issues/38401
class C {
    constructor(/** @public */ x, /** @protected */ y, /** @private */ z) {
    }
}


//// [jsdocAccessibilityTagDeclarations.js]
class Protected {
    constructor(c) {
        this.c = c;
    }
    m() {
        return this.c;
    }
    get p() { return this.c; }
    set p(value) { this.c = value; }
}
class Private {
    constructor(c) {
        this.c = c;
    }
    m() {
        return this.c;
    }
    get p() { return this.c; }
    set p(value) { this.c = value; }
}
class C {
    constructor(x, y, z) {
    }
}
