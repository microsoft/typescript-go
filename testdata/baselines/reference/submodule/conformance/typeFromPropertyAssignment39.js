//// [tests/cases/conformance/salsa/typeFromPropertyAssignment39.ts] ////

//// [a.js]
const foo = {};
foo["baz"] = {};
foo["baz"]["blah"] = 3;




//// [a.d.ts]
const foo: {
    baz: {
        blah: number;
    };
};
declare namespace foo {
    var baz: {
        blah: number;
    };
}
declare namespace foo {
    var blah: number;
}


//// [DtsFileErrors]


a.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
a.d.ts(1,7): error TS2451: Cannot redeclare block-scoped variable 'foo'.
a.d.ts(6,19): error TS2451: Cannot redeclare block-scoped variable 'foo'.
a.d.ts(11,19): error TS2451: Cannot redeclare block-scoped variable 'foo'.


==== a.d.ts (4 errors) ====
    const foo: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
          ~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'foo'.
        baz: {
            blah: number;
        };
    };
    declare namespace foo {
                      ~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'foo'.
        var baz: {
            blah: number;
        };
    }
    declare namespace foo {
                      ~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'foo'.
        var blah: number;
    }
    