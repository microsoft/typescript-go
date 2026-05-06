//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsClassStatic2.ts] ////

//// [Foo.js]
class Base {
  static foo = "";
}
export class Foo extends Base {}
Foo.foo = "foo";

//// [Bar.ts]
import { Foo } from "./Foo.js";

class Bar extends Foo {}
Bar.foo = "foo";




//// [Foo.d.ts]
class Base {
    static foo: string;
}
export class Foo extends Base {
}
export declare namespace Foo {
    var foo: string;
}
//// [Bar.d.ts]
export {};


//// [DtsFileErrors]


out/Foo.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/Foo.d.ts (1 errors) ====
    class Base {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        static foo: string;
    }
    export class Foo extends Base {
    }
    export declare namespace Foo {
        var foo: string;
    }
    
==== out/Bar.d.ts (0 errors) ====
    export {};
    