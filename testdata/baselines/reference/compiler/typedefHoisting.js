//// [tests/cases/compiler/typedefHoisting.ts] ////

//// [x.js]
class C {
    /** @typedef {string} Foo */
    /** @type {Foo} */
    foo = "abc"
}




//// [x.d.ts]
type Foo = string;
declare class C {
    /** @typedef {string} Foo */
    /** @type {Foo} */
    foo: Foo;
}
