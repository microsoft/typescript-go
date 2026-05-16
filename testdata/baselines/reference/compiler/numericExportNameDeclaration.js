//// [tests/cases/compiler/numericExportNameDeclaration.ts] ////

//// [bug.js]
exports[1] = 2;
module.exports[1] = 2;
Object.defineProperty(exports, 1, {});




//// [bug.d.ts]
declare const _exported: any;
export { _exported as 1 };
declare const _exported_1: any;
export { _exported_1 as 1 };
declare const _exported_2: any;
export { _exported_2 as 1 };


//// [DtsFileErrors]


bug.d.ts(2,22): error TS2300: Duplicate identifier '(Missing)'.
bug.d.ts(2,23): error TS1003: Identifier expected.
bug.d.ts(2,25): error TS1128: Declaration or statement expected.
bug.d.ts(4,24): error TS2300: Duplicate identifier '(Missing)'.
bug.d.ts(4,25): error TS1003: Identifier expected.
bug.d.ts(4,27): error TS1128: Declaration or statement expected.
bug.d.ts(6,24): error TS2300: Duplicate identifier '(Missing)'.
bug.d.ts(6,25): error TS1003: Identifier expected.
bug.d.ts(6,27): error TS1128: Declaration or statement expected.


==== bug.d.ts (9 errors) ====
    declare const _exported: any;
    export { _exported as 1 };
                         
!!! error TS2300: Duplicate identifier '(Missing)'.
                          ~
!!! error TS1003: Identifier expected.
                            ~
!!! error TS1128: Declaration or statement expected.
    declare const _exported_1: any;
    export { _exported_1 as 1 };
                           
!!! error TS2300: Duplicate identifier '(Missing)'.
                            ~
!!! error TS1003: Identifier expected.
                              ~
!!! error TS1128: Declaration or statement expected.
    declare const _exported_2: any;
    export { _exported_2 as 1 };
                           
!!! error TS2300: Duplicate identifier '(Missing)'.
                            ~
!!! error TS1003: Identifier expected.
                              ~
!!! error TS1128: Declaration or statement expected.
    