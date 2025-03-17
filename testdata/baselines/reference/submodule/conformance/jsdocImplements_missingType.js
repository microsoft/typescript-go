//// [tests/cases/conformance/jsdoc/jsdocImplements_missingType.ts] ////

//// [a.js]
class A { constructor() { this.x = 0; } }
/** @implements */
class B  {
}


//// [a.js]
class A {
    constructor() { this.x = 0; }
}
class B {
}
