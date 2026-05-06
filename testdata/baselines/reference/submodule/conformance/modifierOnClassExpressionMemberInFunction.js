//// [tests/cases/conformance/classes/classExpressions/modifierOnClassExpressionMemberInFunction.ts] ////

//// [modifierOnClassExpressionMemberInFunction.ts]
function g() {
    var x = class C {
        public prop1 = 1;
        private foo() { }
        static prop2 = 43;
    }
}

//// [modifierOnClassExpressionMemberInFunction.js]
"use strict";
function g() {
    var _a;
    var x = (_a = class C {
            constructor() {
                this.prop1 = 1;
            }
            foo() { }
        },
        _a.prop2 = 43,
        _a);
}


//// [modifierOnClassExpressionMemberInFunction.d.ts]
function g(): void;


//// [DtsFileErrors]


modifierOnClassExpressionMemberInFunction.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== modifierOnClassExpressionMemberInFunction.d.ts (1 errors) ====
    function g(): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    