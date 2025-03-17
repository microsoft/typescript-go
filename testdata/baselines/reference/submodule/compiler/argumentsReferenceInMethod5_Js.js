//// [tests/cases/compiler/argumentsReferenceInMethod5_Js.ts] ////

//// [a.js]
const bar = {
	arguments: {}
}

class A {
	/**
	 * @param {object} [foo={}]
	 */
	m(foo = {}) {
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
    m(foo = {}) {
        this.foo = foo;
        this.bar = bar.arguments;
    }
}
