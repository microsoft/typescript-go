//// [tests/cases/compiler/inferredNonidentifierTypesGetQuotes.ts] ////

//// [inferredNonidentifierTypesGetQuotes.ts]
var x = [{ "a-b": "string" }, {}];

var y = [{ ["a-b"]: "string" }, {}];

//// [inferredNonidentifierTypesGetQuotes.js]
"use strict";
var x = [{ "a-b": "string" }, {}];
var y = [{ ["a-b"]: "string" }, {}];


//// [inferredNonidentifierTypesGetQuotes.d.ts]
var x: ({
    "a-b": string;
} | {
    "a-b"?: undefined;
})[];
var y: ({
    "a-b": string;
} | {
    "a-b"?: undefined;
})[];


//// [DtsFileErrors]


inferredNonidentifierTypesGetQuotes.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== inferredNonidentifierTypesGetQuotes.d.ts (1 errors) ====
    var x: ({
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        "a-b": string;
    } | {
        "a-b"?: undefined;
    })[];
    var y: ({
        "a-b": string;
    } | {
        "a-b"?: undefined;
    })[];
    