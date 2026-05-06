//// [tests/cases/compiler/privacyCheckTypeOfFunction.ts] ////

//// [privacyCheckTypeOfFunction.ts]
function foo() {
}
export var x: typeof foo;
export var b = foo;


//// [privacyCheckTypeOfFunction.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.b = exports.x = void 0;
function foo() {
}
exports.b = foo;


//// [privacyCheckTypeOfFunction.d.ts]
function foo(): void;
export var x: typeof foo;
export var b: typeof foo;
export {};


//// [DtsFileErrors]


privacyCheckTypeOfFunction.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== privacyCheckTypeOfFunction.d.ts (1 errors) ====
    function foo(): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export var x: typeof foo;
    export var b: typeof foo;
    export {};
    