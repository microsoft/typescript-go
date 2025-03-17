//// [tests/cases/compiler/jsDeclarationsInheritedTypes.ts] ////

//// [a.js]
/**
 * @typedef A
 * @property {string} a
 */

/**
 * @typedef B
 * @property {number} b
 */

 class C1 {
    /**
     * @type {A}
     */
    value;
}

class C2 extends C1 {
    /**
     * @type {A}
     */
    value;
}

class C3 extends C1 {
    /**
     * @type {A & B}
     */
    value;
}


//// [a.js]
class C1 {
    value;
}
class C2 extends C1 {
    value;
}
class C3 extends C1 {
    value;
}
