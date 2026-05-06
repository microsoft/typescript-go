//// [tests/cases/conformance/salsa/privateIdentifierExpando.ts] ////

//// [privateIdentifierExpando.js]
const x = {};
x.#bar.baz = 20;




//// [privateIdentifierExpando.d.ts]
const x: {};


//// [DtsFileErrors]


privateIdentifierExpando.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== privateIdentifierExpando.d.ts (1 errors) ====
    const x: {};
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    