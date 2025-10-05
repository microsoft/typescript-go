//// [tests/cases/compiler/isolatedDeclarationErrors.ts] ////

//// [isolatedDeclarationErrors.ts]
function errorOnAssignmentBelowDecl(): void {}
errorOnAssignmentBelowDecl.a = "";

const errorOnAssignmentBelow = (): void => {}
errorOnAssignmentBelow.a = "";

const errorOnMissingReturn = () => {}
errorOnMissingReturn.a = "";


//// [isolatedDeclarationErrors.js]
function errorOnAssignmentBelowDecl() { }
errorOnAssignmentBelowDecl.a = "";
const errorOnAssignmentBelow = () => { };
errorOnAssignmentBelow.a = "";
const errorOnMissingReturn = () => { };
errorOnMissingReturn.a = "";


//// [isolatedDeclarationErrors.d.ts]
declare function errorOnAssignmentBelowDecl(): void;
declare namespace errorOnAssignmentBelowDecl {
    const a: "";
}
declare const errorOnAssignmentBelow: {
    (): void;
    a: string;
};
declare namespace errorOnAssignmentBelow {
    const a: "";
}
declare const errorOnMissingReturn: {
    (): void;
    a: string;
};
declare namespace errorOnMissingReturn {
    const a: "";
}


//// [DtsFileErrors]


isolatedDeclarationErrors.d.ts(5,15): error TS2451: Cannot redeclare block-scoped variable 'errorOnAssignmentBelow'.
isolatedDeclarationErrors.d.ts(9,19): error TS2451: Cannot redeclare block-scoped variable 'errorOnAssignmentBelow'.
isolatedDeclarationErrors.d.ts(12,15): error TS2451: Cannot redeclare block-scoped variable 'errorOnMissingReturn'.
isolatedDeclarationErrors.d.ts(16,19): error TS2451: Cannot redeclare block-scoped variable 'errorOnMissingReturn'.


==== isolatedDeclarationErrors.d.ts (4 errors) ====
    declare function errorOnAssignmentBelowDecl(): void;
    declare namespace errorOnAssignmentBelowDecl {
        const a: "";
    }
    declare const errorOnAssignmentBelow: {
                  ~~~~~~~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'errorOnAssignmentBelow'.
        (): void;
        a: string;
    };
    declare namespace errorOnAssignmentBelow {
                      ~~~~~~~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'errorOnAssignmentBelow'.
        const a: "";
    }
    declare const errorOnMissingReturn: {
                  ~~~~~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'errorOnMissingReturn'.
        (): void;
        a: string;
    };
    declare namespace errorOnMissingReturn {
                      ~~~~~~~~~~~~~~~~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'errorOnMissingReturn'.
        const a: "";
    }
    