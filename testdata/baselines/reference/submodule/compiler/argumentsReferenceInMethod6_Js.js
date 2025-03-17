//// [tests/cases/compiler/argumentsReferenceInMethod6_Js.ts] ////

//// [a.js]
class A {
	m() {
		/**
		 * @type object
		 */
		this.foo = arguments;
	}
}


//// [a.js]
class A {
    m() {
        this.foo = arguments;
    }
}
