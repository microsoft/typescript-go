//// [tests/cases/compiler/objectLiteralDeclarationGeneration1.ts] ////

//// [objectLiteralDeclarationGeneration1.ts]
class y<T extends {}>{ }

//// [objectLiteralDeclarationGeneration1.js]
"use strict";
class y {
}


//// [objectLiteralDeclarationGeneration1.d.ts]
class y<T extends {}> {
}


//// [DtsFileErrors]


objectLiteralDeclarationGeneration1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== objectLiteralDeclarationGeneration1.d.ts (1 errors) ====
    class y<T extends {}> {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    