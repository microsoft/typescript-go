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




//// [a.d.ts]
class A {
    /**
     * @type object
     */
    arguments: object;
    /**
     * @param {object} [foo={}]
     */
    m(foo?: object): void;
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
         * @param {object} [foo={}]
         */
        m(foo?: object): void;
    }
    