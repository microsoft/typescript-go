//// [tests/cases/conformance/types/namedTypes/optionalMethods.ts] ////

//// [optionalMethods.ts]
interface Foo {
    a: number;
    b?: number;
    f(): number;
    g?(): number;
}

function test1(x: Foo) {
    x.a;
    x.b;
    x.f;
    x.g;
    let f1 = x.f();
    let g1 = x.g && x.g();
    let g2 = x.g ? x.g() : 0;
}

class Bar {
    a: number = 0;
    b?: number;
    c? = 2;
    constructor(public d?: number, public e = 10) {}
    f() {
        return 1;
    }
    g?(): number;  // Body of optional method can be omitted
    h?() {
        return 2;
    }
}

function test2(x: Bar) {
    x.a;
    x.b;
    x.c;
    x.d;
    x.e;
    x.f;
    x.g;
    let f1 = x.f();
    let g1 = x.g && x.g();
    let g2 = x.g ? x.g() : 0;
    let h1 = x.h && x.h();
    let h2 = x.h ? x.h() : 0;
}

class Base {
    a?: number;
    f?(): number;
}

class Derived extends Base {
    a = 1;
    f(): number { return 1; }
}


//// [optionalMethods.js]
"use strict";
function test1(x) {
    x.a;
    x.b;
    x.f;
    x.g;
    let f1 = x.f();
    let g1 = x.g && x.g();
    let g2 = x.g ? x.g() : 0;
}
class Bar {
    constructor(d, e = 10) {
        this.d = d;
        this.e = e;
        this.a = 0;
        this.c = 2;
    }
    f() {
        return 1;
    }
    h() {
        return 2;
    }
}
function test2(x) {
    x.a;
    x.b;
    x.c;
    x.d;
    x.e;
    x.f;
    x.g;
    let f1 = x.f();
    let g1 = x.g && x.g();
    let g2 = x.g ? x.g() : 0;
    let h1 = x.h && x.h();
    let h2 = x.h ? x.h() : 0;
}
class Base {
}
class Derived extends Base {
    constructor() {
        super(...arguments);
        this.a = 1;
    }
    f() { return 1; }
}


//// [optionalMethods.d.ts]
interface Foo {
    a: number;
    b?: number;
    f(): number;
    g?(): number;
}
function test1(x: Foo): void;
class Bar {
    d?: number | undefined;
    e: number;
    a: number;
    b?: number;
    c?: number | undefined;
    constructor(d?: number | undefined, e?: number);
    f(): number;
    g?(): number;
    h?(): number;
}
function test2(x: Bar): void;
class Base {
    a?: number;
    f?(): number;
}
class Derived extends Base {
    a: number;
    f(): number;
}


//// [DtsFileErrors]


optionalMethods.d.ts(7,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== optionalMethods.d.ts (1 errors) ====
    interface Foo {
        a: number;
        b?: number;
        f(): number;
        g?(): number;
    }
    function test1(x: Foo): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    class Bar {
        d?: number | undefined;
        e: number;
        a: number;
        b?: number;
        c?: number | undefined;
        constructor(d?: number | undefined, e?: number);
        f(): number;
        g?(): number;
        h?(): number;
    }
    function test2(x: Bar): void;
    class Base {
        a?: number;
        f?(): number;
    }
    class Derived extends Base {
        a: number;
        f(): number;
    }
    