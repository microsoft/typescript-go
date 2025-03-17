//// [tests/cases/compiler/declarationEmitInvalidReferenceAllowJs.ts] ////

//// [declarationEmitInvalidReferenceAllowJs.ts]
/// <reference path="invalid" />
var x = 0; 


//// [declarationEmitInvalidReferenceAllowJs.js]
var x = 0;
//// [invalid.js]
