//// [tests/cases/conformance/jsdoc/jsdocReadonlyDeclarations.ts] ////

//// [jsdocReadonlyDeclarations.js]
class C {
    /** @readonly */
    x = 6
    /** @readonly */
    constructor(n) {
        this.x = n
        /**
         * @readonly
         * @type {number}
         */
        this.y = n
    }
}
new C().x

function F() {
    /** @readonly */
    this.z = 1
}

// https://github.com/microsoft/TypeScript/issues/38401
class D {
    constructor(/** @readonly */ x) {}
}


//// [jsdocReadonlyDeclarations.js]
class C {
    x = 6;
    constructor(n) {
        this.x = n;
        this.y = n;
    }
}
new C().x;
function F() {
    this.z = 1;
}
class D {
    constructor(x) { }
}
