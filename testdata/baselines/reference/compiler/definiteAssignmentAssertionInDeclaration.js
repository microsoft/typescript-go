//// [tests/cases/compiler/definiteAssignmentAssertionInDeclaration.ts] ////

//// [definiteAssignmentAssertionInDeclaration.ts]
export class DbObject {
    id!: string;
}

//// [definiteAssignmentAssertionInDeclaration.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DbObject = void 0;
class DbObject {
    id;
}
exports.DbObject = DbObject;


//// [definiteAssignmentAssertionInDeclaration.d.ts]
export declare class DbObject {
    id: string;
}
