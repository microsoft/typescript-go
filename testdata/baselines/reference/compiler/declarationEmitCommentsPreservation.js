//// [tests/cases/compiler/declarationEmitCommentsPreservation.ts] ////

//// [declarationEmitCommentsPreservation.ts]
// Comment
export class DbObject {
    // Comment
    id: string = ""; // Comment
    // Comment
    method() { }
}

//// [declarationEmitCommentsPreservation.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DbObject = void 0;
// Comment
class DbObject {
    // Comment
    id = ""; // Comment
    // Comment
    method() { }
}
exports.DbObject = DbObject;


//// [declarationEmitCommentsPreservation.d.ts]
export declare class DbObject {
    id: string;
    method(): void;
}
