//// [tests/cases/compiler/moduleOuterQualification.ts] ////

//// [moduleOuterQualification.ts]
declare namespace outer {
  interface Beta { }
  namespace inner {
    // .d.ts emit: should be 'extends outer.Beta'
    export interface Beta extends outer.Beta { }
  }
}


//// [moduleOuterQualification.js]
"use strict";


//// [moduleOuterQualification.d.ts]
namespace outer {
    interface Beta {
    }
    namespace inner {
        interface Beta extends outer.Beta {
        }
    }
}


//// [DtsFileErrors]


moduleOuterQualification.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== moduleOuterQualification.d.ts (1 errors) ====
    namespace outer {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        interface Beta {
        }
        namespace inner {
            interface Beta extends outer.Beta {
            }
        }
    }
    