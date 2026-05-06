//// [tests/cases/compiler/declarationEmitPrivateSymbolCausesVarDeclarationToBeEmitted.ts] ////

//// [declarationEmitPrivateSymbolCausesVarDeclarationToBeEmitted.ts]
const _data = Symbol('data');

export class User {
    private [_data] : any;
};


//// [declarationEmitPrivateSymbolCausesVarDeclarationToBeEmitted.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.User = void 0;
const _data = Symbol('data');
class User {
}
exports.User = User;
;


//// [declarationEmitPrivateSymbolCausesVarDeclarationToBeEmitted.d.ts]
const _data: unique symbol;
export class User {
    private [_data];
}
export {};


//// [DtsFileErrors]


declarationEmitPrivateSymbolCausesVarDeclarationToBeEmitted.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitPrivateSymbolCausesVarDeclarationToBeEmitted.d.ts (1 errors) ====
    const _data: unique symbol;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export class User {
        private [_data];
    }
    export {};
    