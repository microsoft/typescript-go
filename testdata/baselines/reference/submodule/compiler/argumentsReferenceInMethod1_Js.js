//// [tests/cases/compiler/argumentsReferenceInMethod1_Js.ts] ////

//// [a.js]
class A {
	/**
	 * @param {object} [foo={}]
	 */
	m(foo = {}) {
		/**
		 * @type object
		 */
		this.arguments = foo;
	}
}


//// [a.js]
class A {
    m(foo = {}) {
        this.arguments = foo;
    }
}
