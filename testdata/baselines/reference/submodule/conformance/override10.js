//// [tests/cases/conformance/override/override10.ts] ////

//// [override10.ts]
abstract class Base {
    abstract foo(): unknown;
    abstract bar(): void;
}

abstract class Sub extends Base {
    abstract override foo(): number;
    bar() { }
}

//// [override10.js]
"use strict";
class Base {
}
class Sub extends Base {
    bar() { }
}


//// [override10.d.ts]
abstract class Base {
    abstract foo(): unknown;
    abstract bar(): void;
}
abstract class Sub extends Base {
    abstract foo(): number;
    bar(): void;
}


//// [DtsFileErrors]


override10.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== override10.d.ts (1 errors) ====
    abstract class Base {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        abstract foo(): unknown;
        abstract bar(): void;
    }
    abstract class Sub extends Base {
        abstract foo(): number;
        bar(): void;
    }
    