//// [tests/cases/compiler/declFilePrivateStatic.ts] ////

//// [declFilePrivateStatic.ts]
class C {
    private static x = 1;
    static y = 1;

    private static a() { }
    static b() { }

    private static get c() { return 1; }
    static get d() { return 1; }

    private static set e(v) { }
    static set f(v) { }
}

//// [declFilePrivateStatic.js]
"use strict";
class C {
    static a() { }
    static b() { }
    static get c() { return 1; }
    static get d() { return 1; }
    static set e(v) { }
    static set f(v) { }
}
C.x = 1;
C.y = 1;


//// [declFilePrivateStatic.d.ts]
class C {
    private static x;
    static y: number;
    private static a;
    static b(): void;
    private static get c();
    static get d(): number;
    private static set e(value);
    static set f(v: any);
}


//// [DtsFileErrors]


declFilePrivateStatic.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFilePrivateStatic.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        private static x;
        static y: number;
        private static a;
        static b(): void;
        private static get c();
        static get d(): number;
        private static set e(value);
        static set f(v: any);
    }
    