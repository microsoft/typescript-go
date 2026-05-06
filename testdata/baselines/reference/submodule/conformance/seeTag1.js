//// [tests/cases/conformance/jsdoc/seeTag1.ts] ////

//// [seeTag1.ts]
interface Foo {
    foo: string
}

namespace NS {
    export interface Bar {
        baz: Foo
    }
}

/** @see {Foo} foooo*/
const a = ""

/** @see {NS.Bar} ns.bar*/
const b = ""

/** @see {b} b */
const c = ""


//// [seeTag1.js]
"use strict";
/** @see {Foo} foooo*/
const a = "";
/** @see {NS.Bar} ns.bar*/
const b = "";
/** @see {b} b */
const c = "";


//// [seeTag1.d.ts]
interface Foo {
    foo: string;
}
namespace NS {
    interface Bar {
        baz: Foo;
    }
}
/** @see {Foo} foooo*/
const a = "";
/** @see {NS.Bar} ns.bar*/
const b = "";
/** @see {b} b */
const c = "";


//// [DtsFileErrors]


seeTag1.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== seeTag1.d.ts (1 errors) ====
    interface Foo {
        foo: string;
    }
    namespace NS {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        interface Bar {
            baz: Foo;
        }
    }
    /** @see {Foo} foooo*/
    const a = "";
    /** @see {NS.Bar} ns.bar*/
    const b = "";
    /** @see {b} b */
    const c = "";
    