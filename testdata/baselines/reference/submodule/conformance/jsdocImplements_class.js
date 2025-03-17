//// [tests/cases/conformance/jsdoc/jsdocImplements_class.ts] ////

//// [a.js]
class A {
    /** @return {number} */
    method() { throw new Error(); }
}
/** @implements {A} */
class B  {
    method() { return 0 }
}

/** @implements A */
class B2  {
    /** @return {string} */
    method() { return "" }
}

/** @implements {A} */
class B3  {
}


var Ns = {};
/** @implements {A} */
Ns.C1 = class {
    method() { return 11; }
}
/** @implements {A} */
var C2 = class {
    method() { return 12; }
}
var o = {
    /** @implements {A} */
    C3: class {
        method() { return 13; }
    }
}
class CC {
    /** @implements {A} */
    C4 = class {
        method() {
            return 14;
        }
    }
}

var C5;
/** @implements {A} */
Ns.C5 = C5 || class {
    method() {
        return 15;
    }
}


//// [a.js]
class A {
    method() { throw new Error(); }
}
class B {
    method() { return 0; }
}
class B2 {
    method() { return ""; }
}
class B3 {
}
var Ns = {};
Ns.C1 = class {
    method() { return 11; }
};
var C2 = class {
    method() { return 12; }
};
var o = {
    C3: class {
        method() { return 13; }
    }
};
class CC {
    C4 = class {
        method() {
            return 14;
        }
    };
}
var C5;
Ns.C5 = C5 || class {
    method() {
        return 15;
    }
};
