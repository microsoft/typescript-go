//// [tests/cases/compiler/properties.ts] ////

//// [properties.ts]
class MyClass
{
    public get Count(): number
    {
        return 42;
    }

    public set Count(value: number)
    {
        //
    }
}

//// [properties.js]
"use strict";
class MyClass {
    get Count() {
        return 42;
    }
    set Count(value) {
        //
    }
}
//# sourceMappingURL=properties.js.map

//// [properties.d.ts]
class MyClass {
    get Count(): number;
    set Count(value: number);
}


//// [DtsFileErrors]


properties.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== properties.d.ts (1 errors) ====
    class MyClass {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        get Count(): number;
        set Count(value: number);
    }
    