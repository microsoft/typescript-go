//// [tests/cases/compiler/stripInternal1.ts] ////

//// [stripInternal1.ts]
class C {
  foo(): void { }
  // @internal
  bar(): void { }
}

//// [stripInternal1.js]
"use strict";
class C {
    foo() { }
    // @internal
    bar() { }
}


//// [stripInternal1.d.ts]
class C {
    foo(): void;
}


//// [DtsFileErrors]


stripInternal1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== stripInternal1.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        foo(): void;
    }
    