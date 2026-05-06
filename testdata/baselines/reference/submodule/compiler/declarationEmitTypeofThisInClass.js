//// [tests/cases/compiler/declarationEmitTypeofThisInClass.ts] ////

//// [declarationEmitTypeofThisInClass.ts]
class Foo {
    public foo!: string
    public bar!: typeof this.foo //Public property 'bar' of exported class has or is using private name 'this'.(4031)
}

//// [declarationEmitTypeofThisInClass.js]
"use strict";
class Foo {
}


//// [declarationEmitTypeofThisInClass.d.ts]
class Foo {
    foo: string;
    bar: typeof this.foo;
}


//// [DtsFileErrors]


declarationEmitTypeofThisInClass.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitTypeofThisInClass.d.ts (1 errors) ====
    class Foo {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        foo: string;
        bar: typeof this.foo;
    }
    