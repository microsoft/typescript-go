//// [tests/cases/compiler/argumentsReferenceInConstructor5_Js.ts] ////

//// [a.js]
const bar = {
	arguments: {}
}

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
		this.foo = foo;

		/**
		 * @type object
		 */
		this.bar = bar.arguments;
	}
}


//// [a.js]
const bar = {
    arguments: {}
};
class A {
    constructor(foo = {}) {
        this.foo = foo;
        this.bar = bar.arguments;
    }
}
