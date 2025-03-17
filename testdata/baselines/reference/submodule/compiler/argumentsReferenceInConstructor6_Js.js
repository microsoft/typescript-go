//// [tests/cases/compiler/argumentsReferenceInConstructor6_Js.ts] ////

//// [a.js]
class A {
	constructor() {
		/**
		 * @type object
		 */
		this.foo = arguments;
	}
}


//// [a.js]
class A {
    constructor() {
        this.foo = arguments;
    }
}
