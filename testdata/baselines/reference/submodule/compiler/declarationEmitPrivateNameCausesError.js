//// [tests/cases/compiler/declarationEmitPrivateNameCausesError.ts] ////

//// [file.ts]
const IGNORE_EXTRA_VARIABLES = Symbol(); //Notice how this is unexported

//This is exported
export function ignoreExtraVariables<CtorT extends {new(...args:any[]):{}}> (ctor : CtorT) {
    return class extends ctor {
        [IGNORE_EXTRA_VARIABLES] = true; //An unexported constant is used
    };
}

//// [file.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ignoreExtraVariables = ignoreExtraVariables;
const IGNORE_EXTRA_VARIABLES = Symbol(); //Notice how this is unexported
//This is exported
function ignoreExtraVariables(ctor) {
    var _a, _b;
    return _b = class extends ctor {
            constructor() {
                super(...arguments);
                this[_a] = true; //An unexported constant is used
            }
        },
        _a = IGNORE_EXTRA_VARIABLES,
        _b;
}


//// [file.d.ts]
const IGNORE_EXTRA_VARIABLES: unique symbol;
export function ignoreExtraVariables<CtorT extends {
    new (...args: any[]): {};
}>(ctor: CtorT): {
    new (...args: any[]): {
        [IGNORE_EXTRA_VARIABLES]: boolean;
    };
} & CtorT;
export {};


//// [DtsFileErrors]


file.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== file.d.ts (1 errors) ====
    const IGNORE_EXTRA_VARIABLES: unique symbol;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export function ignoreExtraVariables<CtorT extends {
        new (...args: any[]): {};
    }>(ctor: CtorT): {
        new (...args: any[]): {
            [IGNORE_EXTRA_VARIABLES]: boolean;
        };
    } & CtorT;
    export {};
    