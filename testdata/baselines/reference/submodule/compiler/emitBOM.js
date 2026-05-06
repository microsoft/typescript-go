//// [tests/cases/compiler/emitBOM.ts] ////

//// [emitBOM.ts]
// JS and d.ts output should have a BOM but not the sourcemap
var x;

//// [emitBOM.js]
﻿"use strict";
// JS and d.ts output should have a BOM but not the sourcemap
var x;
//# sourceMappingURL=emitBOM.js.map

//// [emitBOM.d.ts]
﻿var x: any;


//// [DtsFileErrors]


emitBOM.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== emitBOM.d.ts (1 errors) ====
    var x: any;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    