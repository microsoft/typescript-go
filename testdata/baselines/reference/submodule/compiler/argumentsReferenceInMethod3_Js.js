//// [tests/cases/compiler/argumentsReferenceInMethod3_Js.ts] ////

//// [a.js]
class A {
	get arguments() {
		return { bar: {} };
	}
}

class B extends A {
	/**
	 * @param {object} [foo={}]
	 */
	m(foo = {}) {
		/**
		 * @type object
		 */
		this.x = foo;

		/**
		 * @type object
		 */
		this.y = super.arguments.bar;
	}
}


//// [a.js]
class A {
    get arguments() {
        return { bar: {} };
    }
}
class B extends A {
    m(foo = {}) {
        this.x = foo;
        this.y = super.arguments.bar;
    }
}
