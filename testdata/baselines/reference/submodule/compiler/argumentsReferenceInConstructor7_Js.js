//// [tests/cases/compiler/argumentsReferenceInConstructor7_Js.ts] ////

//// [a.js]
class A {
	constructor() {
		/**
		 * @type Function
		 */
		this.callee = arguments.callee;
	}
}


//// [a.js]
class A {
    constructor() {
        this.callee = arguments.callee;
    }
}
