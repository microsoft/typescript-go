//// [tests/cases/compiler/methodSignatureDeclarationEmit1.ts] ////

//// [methodSignatureDeclarationEmit1.ts]
class C {
  public foo(n: number): void;
  public foo(s: string): void;
  public foo(a: any): void {
  }
}

//// [methodSignatureDeclarationEmit1.js]
"use strict";
class C {
    foo(a) {
    }
}


//// [methodSignatureDeclarationEmit1.d.ts]
class C {
    foo(n: number): void;
    foo(s: string): void;
}


//// [DtsFileErrors]


methodSignatureDeclarationEmit1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== methodSignatureDeclarationEmit1.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        foo(n: number): void;
        foo(s: string): void;
    }
    