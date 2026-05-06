//// [tests/cases/compiler/stringLiteralObjectLiteralDeclaration1.ts] ////

//// [stringLiteralObjectLiteralDeclaration1.ts]
namespace m1 {
  export var n = { 'foo bar': 4 };
}


//// [stringLiteralObjectLiteralDeclaration1.js]
"use strict";
var m1;
(function (m1) {
    m1.n = { 'foo bar': 4 };
})(m1 || (m1 = {}));


//// [stringLiteralObjectLiteralDeclaration1.d.ts]
namespace m1 {
    var n: {
        'foo bar': number;
    };
}


//// [DtsFileErrors]


stringLiteralObjectLiteralDeclaration1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== stringLiteralObjectLiteralDeclaration1.d.ts (1 errors) ====
    namespace m1 {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        var n: {
            'foo bar': number;
        };
    }
    