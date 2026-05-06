//// [tests/cases/compiler/nonNullableTypes1.ts] ////

//// [nonNullableTypes1.ts]
function f1<T>(x: T) {
    let y = x || "hello";  // NonNullable<T> | string
}

function error(): never {
    throw new Error();
}

function f2<T>(x: T) {  // NonNullable<T>
    return x || error();
}

function f3(x: unknown) {
    let y = x!;  // {}
}

function f4<T extends { x: string } | undefined>(obj: T) {
    if (obj?.x === "hello") {
        obj;  // NonNullable<T>
    }
    if (obj?.x) {
        obj;  // NonNullable<T>
    }
    if (typeof obj?.x === "string") {
        obj;  // NonNullable<T>
    }
}

class A {
    x = "hello";
    foo() {
        let zz = this?.x;  // string
    }
}


//// [nonNullableTypes1.js]
"use strict";
function f1(x) {
    let y = x || "hello"; // NonNullable<T> | string
}
function error() {
    throw new Error();
}
function f2(x) {
    return x || error();
}
function f3(x) {
    let y = x; // {}
}
function f4(obj) {
    if ((obj === null || obj === void 0 ? void 0 : obj.x) === "hello") {
        obj; // NonNullable<T>
    }
    if (obj === null || obj === void 0 ? void 0 : obj.x) {
        obj; // NonNullable<T>
    }
    if (typeof (obj === null || obj === void 0 ? void 0 : obj.x) === "string") {
        obj; // NonNullable<T>
    }
}
class A {
    constructor() {
        this.x = "hello";
    }
    foo() {
        let zz = this === null || this === void 0 ? void 0 : this.x; // string
    }
}


//// [nonNullableTypes1.d.ts]
function f1<T>(x: T): void;
function error(): never;
function f2<T>(x: T): NonNullable<T>;
function f3(x: unknown): void;
function f4<T extends {
    x: string;
} | undefined>(obj: T): void;
class A {
    x: string;
    foo(): void;
}


//// [DtsFileErrors]


nonNullableTypes1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== nonNullableTypes1.d.ts (1 errors) ====
    function f1<T>(x: T): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function error(): never;
    function f2<T>(x: T): NonNullable<T>;
    function f3(x: unknown): void;
    function f4<T extends {
        x: string;
    } | undefined>(obj: T): void;
    class A {
        x: string;
        foo(): void;
    }
    