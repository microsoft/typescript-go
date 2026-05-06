//// [tests/cases/compiler/isolatedDeclarationErrors.ts] ////

//// [isolatedDeclarationErrors.ts]
function errorOnAssignmentBelowDecl(): void {}
errorOnAssignmentBelowDecl.a = "";

const errorOnAssignmentBelow = (): void => {}
errorOnAssignmentBelow.a = "";

const errorOnMissingReturn = () => {}
errorOnMissingReturn.a = "";


//// [isolatedDeclarationErrors.js]
"use strict";
function errorOnAssignmentBelowDecl() { }
errorOnAssignmentBelowDecl.a = "";
const errorOnAssignmentBelow = () => { };
errorOnAssignmentBelow.a = "";
const errorOnMissingReturn = () => { };
errorOnMissingReturn.a = "";


//// [isolatedDeclarationErrors.d.ts]
function errorOnAssignmentBelowDecl(): void;
declare namespace errorOnAssignmentBelowDecl {
    var a: string;
}
function errorOnAssignmentBelow(): void;
declare namespace errorOnAssignmentBelow {
    var a: string;
}
function errorOnMissingReturn(): void;
declare namespace errorOnMissingReturn {
    var a: string;
}
