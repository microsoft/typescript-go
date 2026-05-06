//// [tests/cases/compiler/argumentsReferenceInConstructor5_Js.ts] ////

//// [a.js]
const bar = {
	arguments: {}
}

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
     * Constructor
     *
     * @param {object} [foo={}]
     */
    constructor(foo?: object);
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
         * Constructor
         *
         * @param {object} [foo={}]
         */
        constructor(foo?: object);
    }
    