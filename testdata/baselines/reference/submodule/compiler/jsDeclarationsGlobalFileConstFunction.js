//// [tests/cases/compiler/jsDeclarationsGlobalFileConstFunction.ts] ////

//// [file.js]
const SomeConstructor = function () {
	this.x = 1;
};

const SomeConstructor2 = function () {
};
SomeConstructor2.staticMember = "str";

const SomeConstructor3 = function () {
	this.x = 1;
};
SomeConstructor3.staticMember = "str";




//// [file.d.ts]
declare const SomeConstructor: () => void;
declare const SomeConstructor2: {
    (): void;
    staticMember: string;
};
declare namespace SomeConstructor2 {
    const staticMember: string;
}
declare const SomeConstructor3: {
    (): void;
    staticMember: string;
};
declare namespace SomeConstructor3 {
    const staticMember: string;
}


//// [DtsFileErrors]


file.d.ts(2,15): error TS2451: Cannot redeclare block-scoped variable 'SomeConstructor2'.
file.d.ts(6,19): error TS2451: Cannot redeclare block-scoped variable 'SomeConstructor2'.
file.d.ts(9,15): error TS2451: Cannot redeclare block-scoped variable 'SomeConstructor3'.
file.d.ts(13,19): error TS2451: Cannot redeclare block-scoped variable 'SomeConstructor3'.


==== file.d.ts (4 errors) ====
    declare const SomeConstructor: () => void;
    declare const SomeConstructor2: {
                  ~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'SomeConstructor2'.
        (): void;
        staticMember: string;
    };
    declare namespace SomeConstructor2 {
                      ~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'SomeConstructor2'.
        const staticMember: string;
    }
    declare const SomeConstructor3: {
                  ~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'SomeConstructor3'.
        (): void;
        staticMember: string;
    };
    declare namespace SomeConstructor3 {
                      ~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'SomeConstructor3'.
        const staticMember: string;
    }
    