//// [tests/cases/conformance/types/keyof/keyofAndForIn.ts] ////

//// [keyofAndForIn.ts]
// Repro from #12513

function f1<K extends string, T>(obj: { [P in K]: T }, k: K) {
    const b = k in obj;
    let k1: K;
    for (k1 in obj) {
        let x1 = obj[k1];
    }
    for (let k2 in obj) {
        let x2 = obj[k2];
    }
}

function f2<T>(obj: { [P in keyof T]: T[P] }, k: keyof T) {
    const b = k in obj;
    let k1: keyof T;
    for (k1 in obj) {
        let x1 = obj[k1];
    }
    for (let k2 in obj) {
        let x2 = obj[k2];
    }
}

function f3<T, K extends keyof T>(obj: { [P in K]: T[P] }, k: K) {
    const b = k in obj;
    let k1: K;
    for (k1 in obj) {
        let x1 = obj[k1];
    }
    for (let k2 in obj) {
        let x2 = obj[k2];
    }
}

//// [keyofAndForIn.js]
"use strict";
// Repro from #12513
function f1(obj, k) {
    const b = k in obj;
    let k1;
    for (k1 in obj) {
        let x1 = obj[k1];
    }
    for (let k2 in obj) {
        let x2 = obj[k2];
    }
}
function f2(obj, k) {
    const b = k in obj;
    let k1;
    for (k1 in obj) {
        let x1 = obj[k1];
    }
    for (let k2 in obj) {
        let x2 = obj[k2];
    }
}
function f3(obj, k) {
    const b = k in obj;
    let k1;
    for (k1 in obj) {
        let x1 = obj[k1];
    }
    for (let k2 in obj) {
        let x2 = obj[k2];
    }
}


//// [keyofAndForIn.d.ts]
function f1<K extends string, T>(obj: {
    [P in K]: T;
}, k: K): void;
function f2<T>(obj: {
    [P in keyof T]: T[P];
}, k: keyof T): void;
function f3<T, K extends keyof T>(obj: {
    [P in K]: T[P];
}, k: K): void;


//// [DtsFileErrors]


keyofAndForIn.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== keyofAndForIn.d.ts (1 errors) ====
    function f1<K extends string, T>(obj: {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        [P in K]: T;
    }, k: K): void;
    function f2<T>(obj: {
        [P in keyof T]: T[P];
    }, k: keyof T): void;
    function f3<T, K extends keyof T>(obj: {
        [P in K]: T[P];
    }, k: K): void;
    