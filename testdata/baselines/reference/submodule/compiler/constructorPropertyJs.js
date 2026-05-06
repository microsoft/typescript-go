//// [tests/cases/compiler/constructorPropertyJs.ts] ////

//// [a.js]
class C {
    /**
     * @param {any} a
     */
    foo(a) {
        this.constructor = a;
    }
}




//// [a.d.ts]
class C {
    /**
     * @param {any} a
     */
    foo(a: any): void;
}


//// [DtsFileErrors]


/a.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== /a.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        /**
         * @param {any} a
         */
        foo(a: any): void;
    }
    