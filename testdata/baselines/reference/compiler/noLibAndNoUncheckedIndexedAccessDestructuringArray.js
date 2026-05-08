//// [tests/cases/compiler/noLibAndNoUncheckedIndexedAccessDestructuringArray.ts] ////

//// [noLibAndNoUncheckedIndexedAccessDestructuringArray.ts]
declare var x: string[];
var [a] = x;


//// [noLibAndNoUncheckedIndexedAccessDestructuringArray.js]
"use strict";
var [a] = x;
