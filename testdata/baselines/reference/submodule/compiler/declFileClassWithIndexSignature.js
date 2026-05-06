//// [tests/cases/compiler/declFileClassWithIndexSignature.ts] ////

//// [declFileClassWithIndexSignature.ts]
class BlockIntrinsics {
    [s: string]: string;
}

//// [declFileClassWithIndexSignature.js]
"use strict";
class BlockIntrinsics {
}


//// [declFileClassWithIndexSignature.d.ts]
class BlockIntrinsics {
    [s: string]: string;
}


//// [DtsFileErrors]


declFileClassWithIndexSignature.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileClassWithIndexSignature.d.ts (1 errors) ====
    class BlockIntrinsics {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        [s: string]: string;
    }
    