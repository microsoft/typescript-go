//// [tests/cases/compiler/argumentsReferenceInConstructor4_Js.ts] ////

//// [a.js]
class A {
	/**
	 * Constructor
	 *
	 * @param {object} [foo={}]
	 */
	constructor(foo = {}) {
		const key = "bar";

		/**
		 * @type object
		 */
		this.foo = foo;

		/**
		 * @type object
		 */
		const arguments = this.arguments;

		/**
		 * @type object
		 */
		this.bar = arguments.bar;

		/**
		 * @type object
		 */
		this.baz = arguments[key];

		/**
		 * @type object
		 */
		this.options = arguments;
	}

	get arguments() {
		return { bar: {} };
	}
}


//// [a.js]
class A {
    constructor(foo = {}) {
        const key = "bar";
        this.foo = foo;
        const arguments = this.arguments;
        this.bar = arguments.bar;
        this.baz = arguments[key];
        this.options = arguments;
    }
    get arguments() {
        return { bar: {} };
    }
}
