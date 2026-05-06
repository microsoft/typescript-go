//// [tests/cases/compiler/spreadExpressionContextualType.ts] ////

//// [spreadExpressionContextualType.ts]
// Repro from #43966

interface Orange {
    name: string;
}

interface Apple {
    name: string;
}

function test<T extends Apple | Orange>(item: T): T {
    return { ...item };
}

function test2<T extends Apple | Orange>(item: T): T {
    const x = { ...item };
    return x;
}


//// [spreadExpressionContextualType.js]
"use strict";
// Repro from #43966
function test(item) {
    return Object.assign({}, item);
}
function test2(item) {
    const x = Object.assign({}, item);
    return x;
}


//// [spreadExpressionContextualType.d.ts]
interface Orange {
    name: string;
}
interface Apple {
    name: string;
}
function test<T extends Apple | Orange>(item: T): T;
function test2<T extends Apple | Orange>(item: T): T;


//// [DtsFileErrors]


spreadExpressionContextualType.d.ts(7,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== spreadExpressionContextualType.d.ts (1 errors) ====
    interface Orange {
        name: string;
    }
    interface Apple {
        name: string;
    }
    function test<T extends Apple | Orange>(item: T): T;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function test2<T extends Apple | Orange>(item: T): T;
    