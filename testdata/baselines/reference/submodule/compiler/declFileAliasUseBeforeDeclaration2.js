//// [tests/cases/compiler/declFileAliasUseBeforeDeclaration2.ts] ////

//// [declFileAliasUseBeforeDeclaration2.ts]
declare module "test" {
    namespace A {
        class C {
        }
    }
    class B extends E {
    }
    import E = A.C;
}

//// [declFileAliasUseBeforeDeclaration2.js]
"use strict";


//// [declFileAliasUseBeforeDeclaration2.d.ts]
module "test" {
    namespace A {
        class C {
        }
    }
    class B extends E {
    }
    import E = A.C;
}


//// [DtsFileErrors]


declFileAliasUseBeforeDeclaration2.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileAliasUseBeforeDeclaration2.d.ts (1 errors) ====
    module "test" {
    ~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        namespace A {
            class C {
            }
        }
        class B extends E {
        }
        import E = A.C;
    }
    