//// [tests/cases/compiler/signatureOverloadsWithComments.ts] ////

//// [signatureOverloadsWithComments.ts]
/**
 * Docs
 */
function Foo() {
    return class Bar {
        /**
         * comment 1
         */
        foo(bar: string): void;
        /**
         * @deprecated This signature is deprecated
         *
         * comment 2
         */
        foo(): string;
        foo(bar?: string): string | void {
            return 'hi'
        }
    }
}




//// [signatureOverloadsWithComments.d.ts]
/**
 * Docs
 */
function Foo(): {
    new (): {
        /**
         * comment 1
         */
        foo(bar: string): void;
        /**
         * @deprecated This signature is deprecated
         *
         * comment 2
         */
        foo(): string;
    };
};


//// [DtsFileErrors]


signatureOverloadsWithComments.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== signatureOverloadsWithComments.d.ts (1 errors) ====
    /**
     * Docs
     */
    function Foo(): {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        new (): {
            /**
             * comment 1
             */
            foo(bar: string): void;
            /**
             * @deprecated This signature is deprecated
             *
             * comment 2
             */
            foo(): string;
        };
    };
    