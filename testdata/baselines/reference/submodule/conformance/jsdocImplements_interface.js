//// [tests/cases/conformance/jsdoc/jsdocImplements_interface.ts] ////

//// [defs.d.ts]
interface A {
    mNumber(): number;
}
//// [a.js]
/** @implements A */
class B {
    mNumber() {
        return 0;
    }
}
/** @implements {A} */
class B2 {
    mNumber() {
        return "";
    }
}
/** @implements A */
class B3 {
}


//// [a.js]
class B {
    mNumber() {
        return 0;
    }
}
class B2 {
    mNumber() {
        return "";
    }
}
class B3 {
}
