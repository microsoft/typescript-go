//// [tests/cases/compiler/out-flag.ts] ////

//// [out-flag.ts]
//// @outFile: bin\

// my class comments
class MyClass
{
    // my function comments
    public Count(): number
    {
        return 42;
    }

    public SetCount(value: number)
    {
        //
    }
}


//// [out-flag.js]
"use strict";
//// @outFile: bin\
// my class comments
class MyClass {
    // my function comments
    Count() {
        return 42;
    }
    SetCount(value) {
        //
    }
}
//# sourceMappingURL=out-flag.js.map

//// [out-flag.d.ts]
class MyClass {
    Count(): number;
    SetCount(value: number): void;
}


//// [DtsFileErrors]


out-flag.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out-flag.d.ts (1 errors) ====
    class MyClass {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        Count(): number;
        SetCount(value: number): void;
    }
    