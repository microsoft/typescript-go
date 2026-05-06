//// [tests/cases/compiler/exportNamespaceDeclarationRetainsVisibility.ts] ////

//// [exportNamespaceDeclarationRetainsVisibility.ts]
namespace X {
    interface A {
        kind: 'a';
    }

    interface B {
        kind: 'b';
    }

    export type C = A | B;
}

export = X;

//// [exportNamespaceDeclarationRetainsVisibility.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });


//// [exportNamespaceDeclarationRetainsVisibility.d.ts]
namespace X {
    interface A {
        kind: 'a';
    }
    interface B {
        kind: 'b';
    }
    export type C = A | B;
    export {};
}
export = X;


//// [DtsFileErrors]


exportNamespaceDeclarationRetainsVisibility.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== exportNamespaceDeclarationRetainsVisibility.d.ts (1 errors) ====
    namespace X {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        interface A {
            kind: 'a';
        }
        interface B {
            kind: 'b';
        }
        export type C = A | B;
        export {};
    }
    export = X;
    