//// [tests/cases/compiler/argumentsReferenceInConstructor3_Js.ts] ////

//// [a.js]
class A {
	get arguments() {
		return { bar: {} };
	}
}

class B extends A {
	/**
	 * Constructor
	 *
	 * @param {object} [foo={}]
	 */
	constructor(foo = {}) {
		super();

		/**
		 * @type object
		 */
		this.foo = foo;

		/**
		 * @type object
		 */
		this.bar = super.arguments.foo;
	}
}


//// [a.js]
class A {
    get arguments() {
        return { bar: {} };
    }
}
class B extends A {
    constructor(foo = {}) {
        super();
        this.foo = foo;
        this.bar = super.arguments.foo;
    }
}
