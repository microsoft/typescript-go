//// [tests/cases/compiler/argumentsReferenceInConstructor3_Js.ts] ////

//// [a.js]
class A {
	get arguments() {
		return { bar: {} };
	}
}

class B extends A {
	/**
	 * Constructor
	 *
	 * @param {object} [foo={}]
	 */
	constructor(foo = {}) {
		super();

		/**
		 * @type object
		 */
		this.foo = foo;

		/**
		 * @type object
		 */
		this.bar = super.arguments.foo;
	}
}




//// [a.d.ts]
class A {
    get arguments(): {
        bar: {};
    };
}
class B extends A {
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
    class A {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        get arguments(): {
            bar: {};
        };
    }
    class B extends A {
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
    