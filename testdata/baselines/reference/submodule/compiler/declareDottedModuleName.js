//// [tests/cases/compiler/declareDottedModuleName.ts] ////

//// [declareDottedModuleName.ts]
namespace M {
    namespace P.Q { } // This shouldnt be emitted
}

namespace M {
    export namespace R.S { }  //This should be emitted
}

namespace T.U { // This needs to be emitted
}

//// [declareDottedModuleName.js]
"use strict";


//// [declareDottedModuleName.d.ts]
namespace M {
}
namespace M {
    namespace R.S { }
}
namespace T.U {
}


//// [DtsFileErrors]


declareDottedModuleName.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declareDottedModuleName.d.ts (1 errors) ====
    namespace M {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    namespace M {
        namespace R.S { }
    }
    namespace T.U {
    }
    