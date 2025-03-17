//// [tests/cases/conformance/jsdoc/extendsTag4.ts] ////

//// [foo.js]
/**
 * @constructor
 */
class A {
    constructor() {}
}

/**
 * @extends {A}
 */


//// [foo.js]
class A {
    constructor() { }
}
