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
    const call: () => void;
    const bind: () => void;
    const caller: () => void;
    const toString: () => void;
    const length: 10;
    const length: 10;
}


//// [DtsFileErrors]


isolatedDeclarationErrorsExpandoFunctions.d.ts(8,11): error TS2451: Cannot redeclare block-scoped variable 'length'.
isolatedDeclarationErrorsExpandoFunctions.d.ts(9,11): error TS2451: Cannot redeclare block-scoped variable 'length'.


==== isolatedDeclarationErrorsExpandoFunctions.d.ts (2 errors) ====
    export declare function foo(): void;
    export declare namespace foo {
        const apply: () => void;
        const call: () => void;
        const bind: () => void;
        const caller: () => void;
        const toString: () => void;
        const length: 10;
              ~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'length'.
        const length: 10;
              ~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'length'.
    }
    