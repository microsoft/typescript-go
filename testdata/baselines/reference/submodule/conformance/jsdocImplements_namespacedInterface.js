//// [tests/cases/conformance/jsdoc/jsdocImplements_namespacedInterface.ts] ////

//// [defs.d.ts]
declare namespace N {
    interface A {
        mNumber(): number;
    }
    interface AT<T> {
        gen(): T;
    }
}
//// [a.js]
/** @implements N.A */
class B {
    mNumber() {
        return 0;
    }
}
/** @implements {N.AT<string>} */
class BAT {
    gen() {
        return "";
    }
}




//// [a.d.ts]
/** @implements N.A */
class B implements N.A {
    mNumber(): number;
}
/** @implements {N.AT<string>} */
class BAT implements N.AT<string> {
    gen(): string;
}


//// [DtsFileErrors]


out/a.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== /defs.d.ts (0 errors) ====
    declare namespace N {
        interface A {
            mNumber(): number;
        }
        interface AT<T> {
            gen(): T;
        }
    }
==== out/a.d.ts (1 errors) ====
    /** @implements N.A */
    class B implements N.A {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        mNumber(): number;
    }
    /** @implements {N.AT<string>} */
    class BAT implements N.AT<string> {
        gen(): string;
    }
    