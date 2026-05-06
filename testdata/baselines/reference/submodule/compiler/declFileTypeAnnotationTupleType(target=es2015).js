//// [tests/cases/compiler/declFileTypeAnnotationTupleType.ts] ////

//// [declFileTypeAnnotationTupleType.ts]
class c {
}
namespace m {
    export class c {
    }
    export class g<T> {
    }
}
class g<T> {
}

// Just the name
var k: [c, m.c] = [new c(), new m.c()];
var l = k;

var x: [g<string>, m.g<number>, () => c] = [new g<string>(), new m.g<number>(), () => new c()];
var y = x;

//// [declFileTypeAnnotationTupleType.js]
"use strict";
class c {
}
var m;
(function (m) {
    class c {
    }
    m.c = c;
    class g {
    }
    m.g = g;
})(m || (m = {}));
class g {
}
// Just the name
var k = [new c(), new m.c()];
var l = k;
var x = [new g(), new m.g(), () => new c()];
var y = x;


//// [declFileTypeAnnotationTupleType.d.ts]
class c {
}
namespace m {
    class c {
    }
    class g<T> {
    }
}
class g<T> {
}
var k: [c, m.c];
var l: [c, m.c];
var x: [g<string>, m.g<number>, () => c];
var y: [g<string>, m.g<number>, () => c];


//// [DtsFileErrors]


declFileTypeAnnotationTupleType.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileTypeAnnotationTupleType.d.ts (1 errors) ====
    class c {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    namespace m {
        class c {
        }
        class g<T> {
        }
    }
    class g<T> {
    }
    var k: [c, m.c];
    var l: [c, m.c];
    var x: [g<string>, m.g<number>, () => c];
    var y: [g<string>, m.g<number>, () => c];
    