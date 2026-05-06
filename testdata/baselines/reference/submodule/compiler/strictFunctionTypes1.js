//// [tests/cases/compiler/strictFunctionTypes1.ts] ////

//// [strictFunctionTypes1.ts]
declare function f1<T>(f1: (x: T) => void, f2: (x: T) => void): (x: T) => void;
declare function f2<T>(obj: T, f1: (x: T) => void, f2: (x: T) => void): T;
declare function f3<T>(obj: T, f1: (x: T) => void, f2: (f: (x: T) => void) => void): T;

interface Func<T> { (x: T): void }

declare function f4<T>(f1: Func<T>, f2: Func<T>): Func<T>;

declare function fo(x: Object): void;
declare function fs(x: string): void;
declare function fx(f: (x: "def") => void): void;

const x1 = f1(fo, fs);  // (x: string) => void
const x2 = f2("abc", fo, fs);  // "abc"
const x3 = f3("abc", fo, fx);  // "abc" | "def"
const x4 = f4(fo, fs);  // Func<string>

declare const never: never;

const x10 = f2(never, fo, fs);  // string
const x11 = f3(never, fo, fx);  // "def"

// Repro from #21112

declare function foo<T>(a: ReadonlyArray<T>): T;
let x = foo([]);  // never

// Modified repros from #26127

interface A { a: string }
interface B extends A { b: string }

declare function acceptUnion(x: A | number): void;
declare function acceptA(x: A): void;

declare let a: A;
declare let b: B;

declare function coAndContra<T>(value: T, func: (t: T) => void): T;

const t1: A = coAndContra(a, acceptUnion);
const t2: B = coAndContra(b, acceptA);
const t3: A = coAndContra(never, acceptA);

declare function coAndContraArray<T>(value: T[], func: (t: T) => void): T[];

const t4: A[] = coAndContraArray([a], acceptUnion);
const t5: B[] = coAndContraArray([b], acceptA);
const t6: A[] = coAndContraArray([], acceptA);


//// [strictFunctionTypes1.js]
"use strict";
const x1 = f1(fo, fs); // (x: string) => void
const x2 = f2("abc", fo, fs); // "abc"
const x3 = f3("abc", fo, fx); // "abc" | "def"
const x4 = f4(fo, fs); // Func<string>
const x10 = f2(never, fo, fs); // string
const x11 = f3(never, fo, fx); // "def"
let x = foo([]); // never
const t1 = coAndContra(a, acceptUnion);
const t2 = coAndContra(b, acceptA);
const t3 = coAndContra(never, acceptA);
const t4 = coAndContraArray([a], acceptUnion);
const t5 = coAndContraArray([b], acceptA);
const t6 = coAndContraArray([], acceptA);


//// [strictFunctionTypes1.d.ts]
function f1<T>(f1: (x: T) => void, f2: (x: T) => void): (x: T) => void;
function f2<T>(obj: T, f1: (x: T) => void, f2: (x: T) => void): T;
function f3<T>(obj: T, f1: (x: T) => void, f2: (f: (x: T) => void) => void): T;
interface Func<T> {
    (x: T): void;
}
function f4<T>(f1: Func<T>, f2: Func<T>): Func<T>;
function fo(x: Object): void;
function fs(x: string): void;
function fx(f: (x: "def") => void): void;
const x1: (x: string) => void;
const x2 = "abc";
const x3: string;
const x4: Func<string>;
const never: never;
const x10: string;
const x11: "def";
function foo<T>(a: ReadonlyArray<T>): T;
let x: never;
interface A {
    a: string;
}
interface B extends A {
    b: string;
}
function acceptUnion(x: A | number): void;
function acceptA(x: A): void;
let a: A;
let b: B;
function coAndContra<T>(value: T, func: (t: T) => void): T;
const t1: A;
const t2: B;
const t3: A;
function coAndContraArray<T>(value: T[], func: (t: T) => void): T[];
const t4: A[];
const t5: B[];
const t6: A[];


//// [DtsFileErrors]


strictFunctionTypes1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== strictFunctionTypes1.d.ts (1 errors) ====
    function f1<T>(f1: (x: T) => void, f2: (x: T) => void): (x: T) => void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function f2<T>(obj: T, f1: (x: T) => void, f2: (x: T) => void): T;
    function f3<T>(obj: T, f1: (x: T) => void, f2: (f: (x: T) => void) => void): T;
    interface Func<T> {
        (x: T): void;
    }
    function f4<T>(f1: Func<T>, f2: Func<T>): Func<T>;
    function fo(x: Object): void;
    function fs(x: string): void;
    function fx(f: (x: "def") => void): void;
    const x1: (x: string) => void;
    const x2 = "abc";
    const x3: string;
    const x4: Func<string>;
    const never: never;
    const x10: string;
    const x11: "def";
    function foo<T>(a: ReadonlyArray<T>): T;
    let x: never;
    interface A {
        a: string;
    }
    interface B extends A {
        b: string;
    }
    function acceptUnion(x: A | number): void;
    function acceptA(x: A): void;
    let a: A;
    let b: B;
    function coAndContra<T>(value: T, func: (t: T) => void): T;
    const t1: A;
    const t2: B;
    const t3: A;
    function coAndContraArray<T>(value: T[], func: (t: T) => void): T[];
    const t4: A[];
    const t5: B[];
    const t6: A[];
    