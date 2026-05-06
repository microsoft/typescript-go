//// [tests/cases/conformance/classes/propertyMemberDeclarations/memberAccessorDeclarations/ambientAccessors.ts] ////

//// [ambientAccessors.ts]
// ok to use accessors in ambient class in ES3
declare class C {
    static get a(): string;
    static set a(value: string);

    private static get b(): string;
    private static set b(foo: string);

    get x(): string;
    set x(value: string);

    private get y(): string;
    private set y(foo: string);
}

//// [ambientAccessors.js]
"use strict";


//// [ambientAccessors.d.ts]
class C {
    static get a(): string;
    static set a(value: string);
    private static get b();
    private static set b(value);
    get x(): string;
    set x(value: string);
    private get y();
    private set y(value);
}


//// [DtsFileErrors]


ambientAccessors.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== ambientAccessors.d.ts (1 errors) ====
    class C {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        static get a(): string;
        static set a(value: string);
        private static get b();
        private static set b(value);
        get x(): string;
        set x(value: string);
        private get y();
        private set y(value);
    }
    