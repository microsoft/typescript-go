//// [tests/cases/conformance/classes/classStaticBlock/classStaticBlock25.ts] ////

//// [classStaticBlock25.ts]
const a = 1;
const b = 2;

class C {
    static {
        const a = 11;

        a;
        b;
    }

    static {
        const a = 11;

        a;
        b;
    }
}


//// [classStaticBlock25.js]
"use strict";
const a = 1;
const b = 2;
class C {
    static {
        const a = 11;
        a;
        b;
    }
    static {
        const a = 11;
        a;
        b;
    }
}
//# sourceMappingURL=classStaticBlock25.js.map

//// [classStaticBlock25.d.ts]
const a = 1;
const b = 2;
class C {
}
//# sourceMappingURL=classStaticBlock25.d.ts.map

//// [DtsFileErrors]


classStaticBlock25.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== classStaticBlock25.d.ts (1 errors) ====
    const a = 1;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    const b = 2;
    class C {
    }
    //# sourceMappingURL=classStaticBlock25.d.ts.map