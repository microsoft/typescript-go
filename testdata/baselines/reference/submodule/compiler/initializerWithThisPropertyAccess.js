//// [tests/cases/compiler/initializerWithThisPropertyAccess.ts] ////

//// [initializerWithThisPropertyAccess.ts]
class A {
    a: number;
    b = this.a;  // Error
    c = () => this.a;
    d = (new A()).a;
    constructor() {
        this.a = 1;
    }
}

class B extends A {
    x = this.a;
}

class C {
    a!: number;
    b = this.a;
}

// Repro from #37979

class Foo {
    private bar: Bar;
    readonly barProp = this.bar.prop;
    constructor() {
        this.bar = new Bar();
    }
}

class Bar {
    readonly prop = false;
}


//// [initializerWithThisPropertyAccess.js]
class A {
    a;
    b = this.a; // Error
    c = () => this.a;
    d = (new A()).a;
    constructor() {
        this.a = 1;
    }
}
class B extends A {
    x = this.a;
}
class C {
    a;
    b = this.a;
}
// Repro from #37979
class Foo {
    bar;
    barProp = this.bar.prop;
    constructor() {
        this.bar = new Bar();
    }
}
class Bar {
    prop = false;
}
