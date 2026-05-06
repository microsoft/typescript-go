//// [tests/cases/compiler/argumentsReferenceInConstructor1_Js.ts] ////

//// [a.js]
class A {
	/**
	 * Constructor
	 *
	 * @param {object} [foo={}]
	 */
	constructor(foo = {}) {
		/**
		 * @type object
		 */
		this.arguments = foo;
	}
}




//// [a.d.ts]
class A {
    /**
     * @type object
     */
    arguments: object;
    /**
     * Constructor
     *
     * @param {object} [foo={}]
     */
    constructor(foo?: object);
}


//// [DtsFileErrors]


/a.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== /a.d.ts (1 errors) ====
    class A {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        /**
         * @type object
         */
        arguments: object;
        /**
         * Constructor
         *
         * @param {object} [foo={}]
         */
        constructor(foo?: object);
    }
    