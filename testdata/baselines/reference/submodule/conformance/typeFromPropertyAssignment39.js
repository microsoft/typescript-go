//// [tests/cases/conformance/salsa/typeFromPropertyAssignment39.ts] ////

//// [a.js]
const foo = {};
foo["baz"] = {};
foo["baz"]["blah"] = 3;




//// [a.d.ts]
declare const foo: {
    baz: {
        blah: number;
    };
};
declare namespace foo {
    const blah: number;
}


//// [DtsFileErrors]


a.d.ts(1,15): error TS2451: Cannot redeclare block-scoped variable 'foo'.
a.d.ts(6,19): error TS2451: Cannot redeclare block-scoped variable 'foo'.


==== a.d.ts (2 errors) ====
    declare const foo: {
                  ~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'foo'.
        baz: {
            blah: number;
        };
    };
    declare namespace foo {
                      ~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'foo'.
        const blah: number;
    }
    