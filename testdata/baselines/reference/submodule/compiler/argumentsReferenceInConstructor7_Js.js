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




//// [a.d.ts]
class A {
    /**
     * @type Function
     */
    callee: Function;
    constructor();
}


//// [DtsFileErrors]


/a.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== /a.d.ts (1 errors) ====
    class A {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        /**
         * @type Function
         */
        callee: Function;
        constructor();
    }
    