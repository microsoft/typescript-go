//// [tests/cases/conformance/es6/computedProperties/computedPropertyNamesDeclarationEmit2_ES5.ts] ////

//// [computedPropertyNamesDeclarationEmit2_ES5.ts]
class C {
    static ["" + ""]() { }
    static get ["" + ""]() { return 0; }
    static set ["" + ""](x) { }
}

//// [computedPropertyNamesDeclarationEmit2_ES5.js]
"use strict";
class C {
    static ["" + ""]() { }
    static get ["" + ""]() { return 0; }
    static set ["" + ""](x) { }
}


//// [computedPropertyNamesDeclarationEmit2_ES5.d.ts]
class C {
}


//// [DtsFileErrors]


computedPropertyNamesDeclarationEmit2_ES5.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== computedPropertyNamesDeclarationEmit2_ES5.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    