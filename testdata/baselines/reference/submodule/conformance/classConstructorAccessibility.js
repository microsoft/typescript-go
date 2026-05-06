//// [tests/cases/conformance/classes/constructorDeclarations/classConstructorAccessibility.ts] ////

//// [classConstructorAccessibility.ts]
class C {
    public constructor(public x: number) { }
}

class D {
    private constructor(public x: number) { }
}

class E {
    protected constructor(public x: number) { }
}

var c = new C(1);
var d = new D(1); // error
var e = new E(1); // error

namespace Generic {
    class C<T> {
        public constructor(public x: T) { }
    }

    class D<T> {
        private constructor(public x: T) { }
    }

    class E<T> {
        protected constructor(public x: T) { }
    }

    var c = new C(1);
    var d = new D(1); // error
    var e = new E(1); // error
}


//// [classConstructorAccessibility.js]
"use strict";
class C {
    constructor(x) {
        this.x = x;
    }
}
class D {
    constructor(x) {
        this.x = x;
    }
}
class E {
    constructor(x) {
        this.x = x;
    }
}
var c = new C(1);
var d = new D(1); // error
var e = new E(1); // error
var Generic;
(function (Generic) {
    class C {
        constructor(x) {
            this.x = x;
        }
    }
    class D {
        constructor(x) {
            this.x = x;
        }
    }
    class E {
        constructor(x) {
            this.x = x;
        }
    }
    var c = new C(1);
    var d = new D(1); // error
    var e = new E(1); // error
})(Generic || (Generic = {}));


//// [classConstructorAccessibility.d.ts]
class C {
    x: number;
    constructor(x: number);
}
class D {
    x: number;
    private constructor();
}
class E {
    x: number;
    protected constructor(x: number);
}
var c: C;
var d: D;
var e: E;
namespace Generic {
}
