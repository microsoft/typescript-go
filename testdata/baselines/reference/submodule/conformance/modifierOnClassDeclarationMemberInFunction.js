//// [tests/cases/conformance/classes/classDeclarations/modifierOnClassDeclarationMemberInFunction.ts] ////

//// [modifierOnClassDeclarationMemberInFunction.ts]
function f() {
    class C {
        public baz = 1;
        static foo() { }
        public bar() { }
    }
}

//// [modifierOnClassDeclarationMemberInFunction.js]
"use strict";
function f() {
    class C {
        constructor() {
            this.baz = 1;
        }
        static foo() { }
        bar() { }
    }
}


//// [modifierOnClassDeclarationMemberInFunction.d.ts]
function f(): void;


//// [DtsFileErrors]


modifierOnClassDeclarationMemberInFunction.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== modifierOnClassDeclarationMemberInFunction.d.ts (1 errors) ====
    function f(): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    