//// [tests/cases/compiler/argumentsReferenceInMethod7_Js.ts] ////

//// [a.js]
class A {
	m() {
		/**
		 * @type Function
		 */
		this.callee = arguments.callee;
	}
}


//// [a.js]
class A {
    m() {
        this.callee = arguments.callee;
    }
}
