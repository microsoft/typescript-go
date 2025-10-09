//// [tests/cases/compiler/isolatedDeclarationErrorsExpandoFunctions.ts] ////

//// [isolatedDeclarationErrorsExpandoFunctions.ts]
export function foo() {}

foo.apply = () => {}
foo.call = ()=> {}
foo.bind = ()=> {}
foo.caller = ()=> {}
foo.toString = ()=> {}
foo.length = 10
foo.length = 10


//// [isolatedDeclarationErrorsExpandoFunctions.js]
export function foo() { }
foo.apply = () => { };
foo.call = () => { };
foo.bind = () => { };
foo.caller = () => { };
foo.toString = () => { };
foo.length = 10;
foo.length = 10;


//// [isolatedDeclarationErrorsExpandoFunctions.d.ts]
export declare function foo(): void;
export declare namespace foo {
    const apply: () => void;
}
export declare namespace foo {
    const call: () => void;
}
export declare namespace foo {
    const bind: () => void;
}
export declare namespace foo {
    const caller: () => void;
}
export declare namespace foo {
    const toString: () => void;
}
export declare namespace foo {
    const length: number;
}
export declare namespace foo {
    const length: number;
}


//// [DtsFileErrors]


isolatedDeclarationErrorsExpandoFunctions.d.ts(18,11): error TS2451: Cannot redeclare block-scoped variable 'length'.
isolatedDeclarationErrorsExpandoFunctions.d.ts(21,11): error TS2451: Cannot redeclare block-scoped variable 'length'.


==== isolatedDeclarationErrorsExpandoFunctions.d.ts (2 errors) ====
    export declare function foo(): void;
    export declare namespace foo {
        const apply: () => void;
    }
    export declare namespace foo {
        const call: () => void;
    }
    export declare namespace foo {
        const bind: () => void;
    }
    export declare namespace foo {
        const caller: () => void;
    }
    export declare namespace foo {
        const toString: () => void;
    }
    export declare namespace foo {
        const length: number;
              ~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'length'.
    }
    export declare namespace foo {
        const length: number;
              ~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'length'.
    }
    