//// [tests/cases/conformance/classes/constructorDeclarations/classConstructorAccessibility4.ts] ////

//// [classConstructorAccessibility4.ts]
class A {
    private constructor() { }

    method() {
        class B {
            method() {
                new A(); // OK
            }
        }

        class C extends A { // OK
        }
    }
}

class D {
    protected constructor() { }

    method() {
        class E {
            method() {
                new D(); // OK
            }
        }

        class F extends D { // OK
        }
    }
}

//// [classConstructorAccessibility4.js]
"use strict";
class A {
    constructor() { }
    method() {
        class B {
            method() {
                new A(); // OK
            }
        }
        class C extends A {
        }
    }
}
class D {
    constructor() { }
    method() {
        class E {
            method() {
                new D(); // OK
            }
        }
        class F extends D {
        }
    }
}


//// [classConstructorAccessibility4.d.ts]
class A {
    private constructor();
    method(): void;
}
class D {
    protected constructor();
    method(): void;
}


//// [DtsFileErrors]


classConstructorAccessibility4.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== classConstructorAccessibility4.d.ts (1 errors) ====
    class A {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        private constructor();
        method(): void;
    }
    class D {
        protected constructor();
        method(): void;
    }
    