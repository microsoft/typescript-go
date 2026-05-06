//// [tests/cases/compiler/argumentsReferenceInMethod6_Js.ts] ////

//// [a.js]
class A {
	m() {
		/**
		 * @type object
		 */
		this.foo = arguments;
	}
}




//// [a.d.ts]
class A {
    /**
     * @type object
     */
    foo: object;
    m(): void;
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
        foo: object;
        m(): void;
    }
    