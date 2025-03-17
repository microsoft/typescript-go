//// [tests/cases/compiler/argumentsReferenceInConstructor1_Js.ts] ////

//// [a.js]
class A {
	/**
	 * Constructor
	 *
	 * @param {object} [foo={}]
	 */
	constructor(foo = {}) {
		/**
		 * @type object
		 */
		this.arguments = foo;
	}
}


//// [a.js]
class A {
    constructor(foo = {}) {
        this.arguments = foo;
    }
}
