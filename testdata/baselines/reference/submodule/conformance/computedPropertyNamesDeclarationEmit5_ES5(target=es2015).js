//// [tests/cases/conformance/es6/computedProperties/computedPropertyNamesDeclarationEmit5_ES5.ts] ////

//// [computedPropertyNamesDeclarationEmit5_ES5.ts]
var v = {
    ["" + ""]: 0,
    ["" + ""]() { },
    get ["" + ""]() { return 0; },
    set ["" + ""](x) { }
}

//// [computedPropertyNamesDeclarationEmit5_ES5.js]
"use strict";
var v = {
    ["" + ""]: 0,
    ["" + ""]() { },
    get ["" + ""]() { return 0; },
    set ["" + ""](x) { }
};


//// [computedPropertyNamesDeclarationEmit5_ES5.d.ts]
var v: {
    [x: string]: any;
};


//// [DtsFileErrors]


computedPropertyNamesDeclarationEmit5_ES5.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== computedPropertyNamesDeclarationEmit5_ES5.d.ts (1 errors) ====
    var v: {
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        [x: string]: any;
    };
    