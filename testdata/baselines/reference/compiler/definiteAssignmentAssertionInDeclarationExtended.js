//// [tests/cases/compiler/definiteAssignmentAssertionInDeclarationExtended.ts] ////

//// [definiteAssignmentAssertionInDeclarationExtended.ts]
export class DbObject {
    id!: string;
    name?: string;
    count: number = 0;
    private secret!: string;
    protected value!: number;
    static config!: boolean;
}

export interface IConfig {
    setting?: boolean; 
    optionalSetting?: string;
}

//// [definiteAssignmentAssertionInDeclarationExtended.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DbObject = void 0;
class DbObject {
    id;
    name;
    count = 0;
    secret;
    value;
    static config;
}
exports.DbObject = DbObject;


//// [definiteAssignmentAssertionInDeclarationExtended.d.ts]
export declare class DbObject {
    id: string;
    name?: string;
    count: number;
    private secret;
    protected value: number;
    static config: boolean;
}
export interface IConfig {
    setting?: boolean;
    optionalSetting?: string;
}
