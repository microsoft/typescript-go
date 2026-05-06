//// [tests/cases/compiler/escapedReservedCompilerNamedIdentifier.ts] ////

//// [escapedReservedCompilerNamedIdentifier.ts]
// double underscores
var __proto__ = 10;
var o = {
    "__proto__": 0
};
var b = o["__proto__"];
var o1 = {
    __proto__: 0
};
var b1 = o1["__proto__"];
// Triple underscores
var ___proto__ = 10;
var o2 = {
    "___proto__": 0
};
var b2 = o2["___proto__"];
var o3 = {
    ___proto__: 0
};
var b3 = o3["___proto__"];
// One underscore
var _proto__ = 10;
var o4 = {
    "_proto__": 0
};
var b4 = o4["_proto__"];
var o5 = {
    _proto__: 0
};
var b5 = o5["_proto__"];

//// [escapedReservedCompilerNamedIdentifier.js]
"use strict";
// double underscores
var __proto__ = 10;
var o = {
    "__proto__": 0
};
var b = o["__proto__"];
var o1 = {
    __proto__: 0
};
var b1 = o1["__proto__"];
// Triple underscores
var ___proto__ = 10;
var o2 = {
    "___proto__": 0
};
var b2 = o2["___proto__"];
var o3 = {
    ___proto__: 0
};
var b3 = o3["___proto__"];
// One underscore
var _proto__ = 10;
var o4 = {
    "_proto__": 0
};
var b4 = o4["_proto__"];
var o5 = {
    _proto__: 0
};
var b5 = o5["_proto__"];


//// [escapedReservedCompilerNamedIdentifier.d.ts]
var __proto__: number;
var o: {
    "__proto__": number;
};
var b: number;
var o1: {
    __proto__: number;
};
var b1: number;
var ___proto__: number;
var o2: {
    "___proto__": number;
};
var b2: number;
var o3: {
    ___proto__: number;
};
var b3: number;
var _proto__: number;
var o4: {
    "_proto__": number;
};
var b4: number;
var o5: {
    _proto__: number;
};
var b5: number;


//// [DtsFileErrors]


escapedReservedCompilerNamedIdentifier.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== escapedReservedCompilerNamedIdentifier.d.ts (1 errors) ====
    var __proto__: number;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    var o: {
        "__proto__": number;
    };
    var b: number;
    var o1: {
        __proto__: number;
    };
    var b1: number;
    var ___proto__: number;
    var o2: {
        "___proto__": number;
    };
    var b2: number;
    var o3: {
        ___proto__: number;
    };
    var b3: number;
    var _proto__: number;
    var o4: {
        "_proto__": number;
    };
    var b4: number;
    var o5: {
        _proto__: number;
    };
    var b5: number;
    