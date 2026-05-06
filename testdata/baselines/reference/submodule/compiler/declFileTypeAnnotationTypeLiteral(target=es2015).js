//// [tests/cases/compiler/declFileTypeAnnotationTypeLiteral.ts] ////

//// [declFileTypeAnnotationTypeLiteral.ts]
class c {
}
class g<T> {
}
namespace m {
    export class c {
    }
}

// Object literal with everything
var x: {
    // Call signatures
    (a: number): c;
    (a: string): g<string>;

    // Construct signatures
    new (a: number): c;
    new (a: string): m.c;

    // Indexers
    [n: number]: c;
    [n: string]: c;

    // Properties
    a: c;
    b: g<string>;

    // methods
    m1(): g<number>;
    m2(a: string, b?: number, ...c: c[]): string;
};


// Function type
var y: (a: string) => string;

// constructor type
var z: new (a: string) => m.c;

//// [declFileTypeAnnotationTypeLiteral.js]
"use strict";
class c {
}
class g {
}
var m;
(function (m) {
    class c {
    }
    m.c = c;
})(m || (m = {}));
// Object literal with everything
var x;
// Function type
var y;
// constructor type
var z;


//// [declFileTypeAnnotationTypeLiteral.d.ts]
class c {
}
class g<T> {
}
namespace m {
    class c {
    }
}
var x: {
    (a: number): c;
    (a: string): g<string>;
    new (a: number): c;
    new (a: string): m.c;
    [n: number]: c;
    [n: string]: c;
    a: c;
    b: g<string>;
    m1(): g<number>;
    m2(a: string, b?: number, ...c: c[]): string;
};
var y: (a: string) => string;
var z: new (a: string) => m.c;


//// [DtsFileErrors]


declFileTypeAnnotationTypeLiteral.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileTypeAnnotationTypeLiteral.d.ts (1 errors) ====
    class c {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    class g<T> {
    }
    namespace m {
        class c {
        }
    }
    var x: {
        (a: number): c;
        (a: string): g<string>;
        new (a: number): c;
        new (a: string): m.c;
        [n: number]: c;
        [n: string]: c;
        a: c;
        b: g<string>;
        m1(): g<number>;
        m2(a: string, b?: number, ...c: c[]): string;
    };
    var y: (a: string) => string;
    var z: new (a: string) => m.c;
    