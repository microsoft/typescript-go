//// [tests/cases/compiler/circularInstantiationExpression.ts] ////

//// [circularInstantiationExpression.ts]
declare function foo<T>(t: T): typeof foo<T>;
foo("");


//// [circularInstantiationExpression.js]
"use strict";
foo("");


//// [circularInstantiationExpression.d.ts]
function foo<T>(t: T): typeof foo<T>;


//// [DtsFileErrors]


circularInstantiationExpression.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== circularInstantiationExpression.d.ts (1 errors) ====
    function foo<T>(t: T): typeof foo<T>;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    