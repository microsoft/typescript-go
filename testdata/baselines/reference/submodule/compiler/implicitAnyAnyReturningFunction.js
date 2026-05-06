//// [tests/cases/compiler/implicitAnyAnyReturningFunction.ts] ////

//// [implicitAnyAnyReturningFunction.ts]
function A() {
    return <any>"";
}

function B() {
    var someLocal: any = {};
    return someLocal;
}

class C {
    public A() {
        return <any>"";
    }

    public B() {
        var someLocal: any = {};
        return someLocal;
    }
}


//// [implicitAnyAnyReturningFunction.js]
"use strict";
function A() {
    return "";
}
function B() {
    var someLocal = {};
    return someLocal;
}
class C {
    A() {
        return "";
    }
    B() {
        var someLocal = {};
        return someLocal;
    }
}


//// [implicitAnyAnyReturningFunction.d.ts]
function A(): any;
function B(): any;
class C {
    A(): any;
    B(): any;
}


//// [DtsFileErrors]


implicitAnyAnyReturningFunction.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== implicitAnyAnyReturningFunction.d.ts (1 errors) ====
    function A(): any;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function B(): any;
    class C {
        A(): any;
        B(): any;
    }
    