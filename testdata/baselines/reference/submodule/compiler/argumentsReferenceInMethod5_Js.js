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




//// [a.d.ts]
const bar: {
    arguments: {};
};
class A {
    /**
     * @type object
     */
    foo: object;
    /**
     * @type object
     */
    bar: object;
    /**
     * @param {object} [foo={}]
     */
    m(foo?: object): void;
}


//// [DtsFileErrors]


/a.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== /a.d.ts (1 errors) ====
    const bar: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        arguments: {};
    };
    class A {
        /**
         * @type object
         */
        foo: object;
        /**
         * @type object
         */
        bar: object;
        /**
         * @param {object} [foo={}]
         */
        m(foo?: object): void;
    }
    