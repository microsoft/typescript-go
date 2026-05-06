//// [tests/cases/conformance/classes/classStaticBlock/classStaticBlockUseBeforeDef4.ts] ////

//// [classStaticBlockUseBeforeDef4.ts]
class C {
    static accessor x;
    static {
        this.x = 1;
    }
    static accessor y = this.x;
    static accessor z;
    static {
        this.z = this.y;
    }
}


//// [classStaticBlockUseBeforeDef4.js]
"use strict";
class C {
    static accessor x;
    static {
        this.x = 1;
    }
    static accessor y = this.x;
    static accessor z;
    static {
        this.z = this.y;
    }
}


//// [classStaticBlockUseBeforeDef4.d.ts]
class C {
    static accessor x: number;
    static accessor y: number;
    static accessor z: number;
}


//// [DtsFileErrors]


classStaticBlockUseBeforeDef4.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== classStaticBlockUseBeforeDef4.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        static accessor x: number;
        static accessor y: number;
        static accessor z: number;
    }
    