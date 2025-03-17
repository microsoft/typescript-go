//// [tests/cases/compiler/doNotEmitTripleSlashCommentsOnNotEmittedNode.ts] ////

//// [file0.ts]
/// <reference path="file1.ts" />
declare var OData: any;

//// [file1.ts]
/// <reference path="file0.ts" />
interface F { }



//// [file1.js]
//// [file0.js]
