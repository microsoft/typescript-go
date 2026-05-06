//// [tests/cases/compiler/jsFileCompilationErrorOnDeclarationsWithJsFileReferenceWithOutDir.ts] ////

//// [a.ts]
class c {
}

//// [b.ts]
/// <reference path="c.js"/>
function foo() {
}

//// [c.js]
function bar() {
}

//// [a.js]
"use strict";
class c {
}
//// [c.js]
"use strict";
function bar() {
}
//// [b.js]
"use strict";
/// <reference path="c.js"/>
function foo() {
}


//// [a.d.ts]
class c {
}
//// [c.d.ts]
function bar(): void;
//// [b.d.ts]
function foo(): void;


//// [DtsFileErrors]


outDir/a.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
outDir/b.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
outDir/c.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== outDir/a.d.ts (1 errors) ====
    class c {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    
==== outDir/b.d.ts (1 errors) ====
    function foo(): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    
==== outDir/c.d.ts (1 errors) ====
    function bar(): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    