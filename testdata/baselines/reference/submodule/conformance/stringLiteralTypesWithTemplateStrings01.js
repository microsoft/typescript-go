//// [tests/cases/conformance/types/stringLiteral/stringLiteralTypesWithTemplateStrings01.ts] ////

//// [stringLiteralTypesWithTemplateStrings01.ts]
let ABC: "ABC" = `ABC`;
let DE_NEWLINE_F: "DE\nF" = `DE
F`;
let G_QUOTE_HI: 'G"HI';
let JK_BACKTICK_L: "JK`L" = `JK\`L`;

//// [stringLiteralTypesWithTemplateStrings01.js]
"use strict";
let ABC = `ABC`;
let DE_NEWLINE_F = `DE
F`;
let G_QUOTE_HI;
let JK_BACKTICK_L = `JK\`L`;


//// [stringLiteralTypesWithTemplateStrings01.d.ts]
let ABC: "ABC";
let DE_NEWLINE_F: "DE\nF";
let G_QUOTE_HI: 'G"HI';
let JK_BACKTICK_L: "JK`L";


//// [DtsFileErrors]


stringLiteralTypesWithTemplateStrings01.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== stringLiteralTypesWithTemplateStrings01.d.ts (1 errors) ====
    let ABC: "ABC";
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    let DE_NEWLINE_F: "DE\nF";
    let G_QUOTE_HI: 'G"HI';
    let JK_BACKTICK_L: "JK`L";
    