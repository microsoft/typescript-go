//// [tests/cases/compiler/computedEnumTypeWidening.ts] ////

//// [computedEnumTypeWidening.ts]
declare function computed(x: number): number;

enum E {
    A = computed(0),
    B = computed(1),
    C = computed(2),
    D = computed(3),
}

function f1() {
    const c1 = E.B;  // Fresh E.B
    let v1 = c1;  // E
    const c2 = c1;  // Fresh E.B
    let v2 = c2;  // E
    const c3: E.B = E.B;  // E.B
    let v3 = c3;  // E.B
    const c4: E.B = c1;  // E.B
    let v4 = c4;  // E.B
}

function f2(cond: boolean) {
    const c1 = cond ? E.A : E.B;  // Fresh E.A | fresh E.B
    const c2: E.A | E.B = c1;  // E.A | E.B
    const c3 = cond ? c1 : c2;  // E.A | E.B
    const c4 = cond ? c3 : E.C;  // E.A | E.B | fresh E.C
    const c5: E.A | E.B | E.C = c4; // E.A | E.B | E.C
    let v1 = c1;  // E
    let v2 = c2;  // E.A | E.B
    let v3 = c3;  // E.A | E.B
    let v4 = c4;  // E
    let v5 = c5;  // E.A | E.B | E.C
}

function f3() {
    const c1 = E.B;
    let v1 = c1;  // E
    const c2: E.B = E.B;
    let v2 = c2;  // E.B
    const c3 = E.B as E.B;
    let v3 = c3;  // E.B
    const c4 = <E.B>E.B;
    let v4 = c4;  // E.B
    const c5 = E.B as const;
    let v5 = c5;  // E.B
}

declare enum E2 { A, B, C, D }

function f4() {
    const c1 = E2.B;  // Fresh E2.B
    let v1 = E.B;  // E2
}

const c1 = E.B;
const c2 = E.B as const;
let v1 = E.B;
let v2 = E.B as const;

class C {
  p1 = E.B;
  p2 = E.B as const;
  readonly p3 = E.B;
  readonly p4 = E.B as const;
}

// Repro from #52531

enum MyEnum { A, B, C }

let val1 = MyEnum.A;
val1 = MyEnum.B;

declare enum MyDeclaredEnum { A, B, C }

let val2 = MyDeclaredEnum.A;
val2 = MyDeclaredEnum.B;


//// [computedEnumTypeWidening.js]
"use strict";
var E;
(function (E) {
    E[E["A"] = computed(0)] = "A";
    E[E["B"] = computed(1)] = "B";
    E[E["C"] = computed(2)] = "C";
    E[E["D"] = computed(3)] = "D";
})(E || (E = {}));
function f1() {
    const c1 = E.B; // Fresh E.B
    let v1 = c1; // E
    const c2 = c1; // Fresh E.B
    let v2 = c2; // E
    const c3 = E.B; // E.B
    let v3 = c3; // E.B
    const c4 = c1; // E.B
    let v4 = c4; // E.B
}
function f2(cond) {
    const c1 = cond ? E.A : E.B; // Fresh E.A | fresh E.B
    const c2 = c1; // E.A | E.B
    const c3 = cond ? c1 : c2; // E.A | E.B
    const c4 = cond ? c3 : E.C; // E.A | E.B | fresh E.C
    const c5 = c4; // E.A | E.B | E.C
    let v1 = c1; // E
    let v2 = c2; // E.A | E.B
    let v3 = c3; // E.A | E.B
    let v4 = c4; // E
    let v5 = c5; // E.A | E.B | E.C
}
function f3() {
    const c1 = E.B;
    let v1 = c1; // E
    const c2 = E.B;
    let v2 = c2; // E.B
    const c3 = E.B;
    let v3 = c3; // E.B
    const c4 = E.B;
    let v4 = c4; // E.B
    const c5 = E.B;
    let v5 = c5; // E.B
}
function f4() {
    const c1 = E2.B; // Fresh E2.B
    let v1 = E.B; // E2
}
const c1 = E.B;
const c2 = E.B;
let v1 = E.B;
let v2 = E.B;
class C {
    constructor() {
        this.p1 = E.B;
        this.p2 = E.B;
        this.p3 = E.B;
        this.p4 = E.B;
    }
}
// Repro from #52531
var MyEnum;
(function (MyEnum) {
    MyEnum[MyEnum["A"] = 0] = "A";
    MyEnum[MyEnum["B"] = 1] = "B";
    MyEnum[MyEnum["C"] = 2] = "C";
})(MyEnum || (MyEnum = {}));
let val1 = MyEnum.A;
val1 = MyEnum.B;
let val2 = MyDeclaredEnum.A;
val2 = MyDeclaredEnum.B;


//// [computedEnumTypeWidening.d.ts]
function computed(x: number): number;
enum E {
    A,
    B,
    C,
    D
}
function f1(): void;
function f2(cond: boolean): void;
function f3(): void;
enum E2 {
    A,
    B,
    C,
    D
}
function f4(): void;
const c1 = E.B;
const c2: E.B;
let v1: E;
let v2: E.B;
class C {
    p1: E;
    p2: E.B;
    readonly p3 = E.B;
    readonly p4: E.B;
}
enum MyEnum {
    A = 0,
    B = 1,
    C = 2
}
let val1: MyEnum;
enum MyDeclaredEnum {
    A,
    B,
    C
}
let val2: MyDeclaredEnum;


//// [DtsFileErrors]


computedEnumTypeWidening.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== computedEnumTypeWidening.d.ts (1 errors) ====
    function computed(x: number): number;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    enum E {
        A,
        B,
        C,
        D
    }
    function f1(): void;
    function f2(cond: boolean): void;
    function f3(): void;
    enum E2 {
        A,
        B,
        C,
        D
    }
    function f4(): void;
    const c1 = E.B;
    const c2: E.B;
    let v1: E;
    let v2: E.B;
    class C {
        p1: E;
        p2: E.B;
        readonly p3 = E.B;
        readonly p4: E.B;
    }
    enum MyEnum {
        A = 0,
        B = 1,
        C = 2
    }
    let val1: MyEnum;
    enum MyDeclaredEnum {
        A,
        B,
        C
    }
    let val2: MyDeclaredEnum;
    