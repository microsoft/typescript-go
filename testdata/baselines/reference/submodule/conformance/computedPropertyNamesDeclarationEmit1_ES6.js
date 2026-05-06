//// [tests/cases/conformance/es6/computedProperties/computedPropertyNamesDeclarationEmit1_ES6.ts] ////

//// [computedPropertyNamesDeclarationEmit1_ES6.ts]
class C {
    ["" + ""]() { }
    get ["" + ""]() { return 0; }
    set ["" + ""](x) { }
}

//// [computedPropertyNamesDeclarationEmit1_ES6.js]
"use strict";
class C {
    ["" + ""]() { }
    get ["" + ""]() { return 0; }
    set ["" + ""](x) { }
}


//// [computedPropertyNamesDeclarationEmit1_ES6.d.ts]
class C {
}


//// [DtsFileErrors]


computedPropertyNamesDeclarationEmit1_ES6.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== computedPropertyNamesDeclarationEmit1_ES6.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    