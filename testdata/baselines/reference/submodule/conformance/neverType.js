//// [tests/cases/conformance/types/never/neverType.ts] ////

//// [neverType.ts]
function error(message: string): never {
    throw new Error(message);
}

function errorVoid(message: string) {
    throw new Error(message);
}

function fail() {
    return error("Something failed");
}

function failOrThrow(shouldFail: boolean) {
    if (shouldFail) {
        return fail();
    }
    throw new Error();
}

function infiniteLoop1() {
    while (true) {
    }
}

function infiniteLoop2(): never {
    while (true) {
    }
}

function move1(direction: "up" | "down") {
    switch (direction) {
        case "up":
            return 1;
        case "down":
            return -1; 
    }
    return error("Should never get here");
}

function move2(direction: "up" | "down") {
    return direction === "up" ? 1 :
        direction === "down" ? -1 :
        error("Should never get here");
}

function check<T>(x: T | undefined) {
    return x || error("Undefined value");
}

class C {
    void1() {
        throw new Error();
    }
    void2() {
        while (true) {}
    }
    never1(): never {
        throw new Error();
    }
    never2(): never {
        while (true) {}
    }
}

function f1(x: string | number) {
    if (typeof x === "boolean") {
        x;  // never
    }
}

function f2(x: string | number) {
    while (true) {
        if (typeof x === "boolean") {
            return x;  // never
        }
    }
}

function test(cb: () => string) {
    let s = cb();
    return s;
}

let errorCallback = () => error("Error callback");

test(() => "hello");
test(() => fail());
test(() => { throw new Error(); })
test(errorCallback);


//// [neverType.js]
"use strict";
function error(message) {
    throw new Error(message);
}
function errorVoid(message) {
    throw new Error(message);
}
function fail() {
    return error("Something failed");
}
function failOrThrow(shouldFail) {
    if (shouldFail) {
        return fail();
    }
    throw new Error();
}
function infiniteLoop1() {
    while (true) {
    }
}
function infiniteLoop2() {
    while (true) {
    }
}
function move1(direction) {
    switch (direction) {
        case "up":
            return 1;
        case "down":
            return -1;
    }
    return error("Should never get here");
}
function move2(direction) {
    return direction === "up" ? 1 :
        direction === "down" ? -1 :
            error("Should never get here");
}
function check(x) {
    return x || error("Undefined value");
}
class C {
    void1() {
        throw new Error();
    }
    void2() {
        while (true) { }
    }
    never1() {
        throw new Error();
    }
    never2() {
        while (true) { }
    }
}
function f1(x) {
    if (typeof x === "boolean") {
        x; // never
    }
}
function f2(x) {
    while (true) {
        if (typeof x === "boolean") {
            return x; // never
        }
    }
}
function test(cb) {
    let s = cb();
    return s;
}
let errorCallback = () => error("Error callback");
test(() => "hello");
test(() => fail());
test(() => { throw new Error(); });
test(errorCallback);


//// [neverType.d.ts]
function error(message: string): never;
function errorVoid(message: string): void;
function fail(): never;
function failOrThrow(shouldFail: boolean): never;
function infiniteLoop1(): void;
function infiniteLoop2(): never;
function move1(direction: "up" | "down"): -1 | 1;
function move2(direction: "up" | "down"): -1 | 1;
function check<T>(x: T | undefined): NonNullable<T>;
class C {
    void1(): void;
    void2(): void;
    never1(): never;
    never2(): never;
}
function f1(x: string | number): void;
function f2(x: string | number): never;
function test(cb: () => string): string;
let errorCallback: () => never;


//// [DtsFileErrors]


neverType.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== neverType.d.ts (1 errors) ====
    function error(message: string): never;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function errorVoid(message: string): void;
    function fail(): never;
    function failOrThrow(shouldFail: boolean): never;
    function infiniteLoop1(): void;
    function infiniteLoop2(): never;
    function move1(direction: "up" | "down"): -1 | 1;
    function move2(direction: "up" | "down"): -1 | 1;
    function check<T>(x: T | undefined): NonNullable<T>;
    class C {
        void1(): void;
        void2(): void;
        never1(): never;
        never2(): never;
    }
    function f1(x: string | number): void;
    function f2(x: string | number): never;
    function test(cb: () => string): string;
    let errorCallback: () => never;
    