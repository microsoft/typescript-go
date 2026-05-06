//// [tests/cases/compiler/declFileForClassWithMultipleBaseClasses.ts] ////

//// [declFileForClassWithMultipleBaseClasses.ts]
class A {
    foo() { }
}

class B {
    bar() { }
}

interface I {
    baz();
}

interface J {
    bat();
}


class D implements I, J {
    baz() { }
    bat() { }
    foo() { }
    bar() { }
}

interface I extends A, B {
}

//// [declFileForClassWithMultipleBaseClasses.js]
"use strict";
class A {
    foo() { }
}
class B {
    bar() { }
}
class D {
    baz() { }
    bat() { }
    foo() { }
    bar() { }
}


//// [declFileForClassWithMultipleBaseClasses.d.ts]
class A {
    foo(): void;
}
class B {
    bar(): void;
}
interface I {
    baz(): any;
}
interface J {
    bat(): any;
}
class D implements I, J {
    baz(): void;
    bat(): void;
    foo(): void;
    bar(): void;
}
interface I extends A, B {
}


//// [DtsFileErrors]


declFileForClassWithMultipleBaseClasses.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileForClassWithMultipleBaseClasses.d.ts (1 errors) ====
    class A {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        foo(): void;
    }
    class B {
        bar(): void;
    }
    interface I {
        baz(): any;
    }
    interface J {
        bat(): any;
    }
    class D implements I, J {
        baz(): void;
        bat(): void;
        foo(): void;
        bar(): void;
    }
    interface I extends A, B {
    }
    