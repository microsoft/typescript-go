//// [tests/cases/compiler/argumentsReferenceInMethod3_Js.ts] ////

//// [a.js]
class A {
	get arguments() {
		return { bar: {} };
	}
}

class B extends A {
	/**
	 * @param {object} [foo={}]
	 */
	m(foo = {}) {
		/**
		 * @type object
		 */
		this.x = foo;

		/**
		 * @type object
		 */
		this.y = super.arguments.bar;
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
    x: object;
    /**
     * @type object
     */
    y: object;
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
        get arguments(): {
            bar: {};
        };
    }
    class B extends A {
        /**
         * @type object
         */
        x: object;
        /**
         * @type object
         */
        y: object;
        /**
         * @param {object} [foo={}]
         */
        m(foo?: object): void;
    }
    