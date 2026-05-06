//// [tests/cases/compiler/namespacesDeclaration1.ts] ////

//// [namespacesDeclaration1.ts]
namespace M {
   export namespace N {
      export namespace M2 {
         export interface I {}
      }
   }
}

//// [namespacesDeclaration1.js]
"use strict";


//// [namespacesDeclaration1.d.ts]
namespace M {
    namespace N {
        namespace M2 {
            interface I {
            }
        }
    }
}


//// [DtsFileErrors]


namespacesDeclaration1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== namespacesDeclaration1.d.ts (1 errors) ====
    namespace M {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        namespace N {
            namespace M2 {
                interface I {
                }
            }
        }
    }
    