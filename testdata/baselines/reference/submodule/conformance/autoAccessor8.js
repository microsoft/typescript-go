//// [tests/cases/conformance/classes/propertyMemberDeclarations/autoAccessor8.ts] ////

//// [autoAccessor8.ts]
class C1 {
    accessor a: any;
    static accessor b: any;
}

declare class C2 {
    accessor a: any;
    static accessor b: any;
}

function f() {
    class C3 {
        accessor a: any;
        static accessor b: any;
    }
    return C3;
}


//// [autoAccessor8.js]
"use strict";
class C1 {
    accessor a;
    static accessor b;
}
function f() {
    class C3 {
        accessor a;
        static accessor b;
    }
    return C3;
}


//// [autoAccessor8.d.ts]
class C1 {
    accessor a: any;
    static accessor b: any;
}
class C2 {
    accessor a: any;
    static accessor b: any;
}
function f(): {
    new (): {
        get a(): any;
        set a(arg: any);
    };
    get b(): any;
    set b(arg: any);
};


//// [DtsFileErrors]


autoAccessor8.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== autoAccessor8.d.ts (1 errors) ====
    class C1 {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        accessor a: any;
        static accessor b: any;
    }
    class C2 {
        accessor a: any;
        static accessor b: any;
    }
    function f(): {
        new (): {
            get a(): any;
            set a(arg: any);
        };
        get b(): any;
        set b(arg: any);
    };
    