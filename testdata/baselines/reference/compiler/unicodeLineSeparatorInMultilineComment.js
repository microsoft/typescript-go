//// [tests/cases/compiler/unicodeLineSeparatorInMultilineComment.ts] ////

//// [unicodeLineSeparatorInMultilineComment.ts]
/* aâ€¨b */ const x = 1;


//// [unicodeLineSeparatorInMultilineComment.js]
"use strict";
/* aâ€
b */ const x = 1;
