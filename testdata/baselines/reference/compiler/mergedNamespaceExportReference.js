//// [tests/cases/compiler/mergedNamespaceExportReference.ts] ////

//// [mergedNamespaceExportReference.ts]
// Test that references to exported namespace members across merged namespace
// declarations are correctly qualified in the emitted JavaScript.

namespace N {
    export function foo() { return 1; }
    export var x = 1;
    export class C {}
}

namespace N {
    // These should emit as N.foo(), N.x, and N.C
    foo();
    x;
    class D extends C {}
}


//// [mergedNamespaceExportReference.js]
// Test that references to exported namespace members across merged namespace
// declarations are correctly qualified in the emitted JavaScript.
var N;
(function (N) {
    function foo() { return 1; }
    N.foo = foo;
    N.x = 1;
    class C {
    }
    N.C = C;
})(N || (N = {}));
(function (N) {
    // These should emit as N.foo(), N.x, and N.C
    N.foo();
    N.x;
    class D extends N.C {
    }
})(N || (N = {}));
