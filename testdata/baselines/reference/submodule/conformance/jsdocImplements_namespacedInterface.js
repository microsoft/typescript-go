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


//// [a.js]
class B {
    mNumber() {
        return 0;
    }
}
class BAT {
    gen() {
        return "";
    }
}
