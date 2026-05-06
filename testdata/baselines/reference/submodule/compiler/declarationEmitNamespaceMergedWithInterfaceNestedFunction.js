//// [tests/cases/compiler/declarationEmitNamespaceMergedWithInterfaceNestedFunction.ts] ////

//// [declarationEmitNamespaceMergedWithInterfaceNestedFunction.ts]
export interface Foo {
    item: Bar;
}

interface Bar {
    baz(): void;
}

namespace Bar {
    export function biz() {
        return 0;
    }
}

//// [declarationEmitNamespaceMergedWithInterfaceNestedFunction.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var Bar;
(function (Bar) {
    function biz() {
        return 0;
    }
    Bar.biz = biz;
})(Bar || (Bar = {}));


//// [declarationEmitNamespaceMergedWithInterfaceNestedFunction.d.ts]
export interface Foo {
    item: Bar;
}
interface Bar {
    baz(): void;
}
namespace Bar {
    function biz(): number;
}
export {};


//// [DtsFileErrors]


declarationEmitNamespaceMergedWithInterfaceNestedFunction.d.ts(7,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitNamespaceMergedWithInterfaceNestedFunction.d.ts (1 errors) ====
    export interface Foo {
        item: Bar;
    }
    interface Bar {
        baz(): void;
    }
    namespace Bar {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        function biz(): number;
    }
    export {};
    