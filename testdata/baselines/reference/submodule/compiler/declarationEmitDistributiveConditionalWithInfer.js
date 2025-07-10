//// [tests/cases/compiler/declarationEmitDistributiveConditionalWithInfer.ts] ////

//// [declarationEmitDistributiveConditionalWithInfer.ts]
// This function's type is changed on declaration
export const fun = (
    subFun: <Collection, Field extends keyof Collection>()
        => FlatArray<Collection[Field], 0>[]) => { };


//// [declarationEmitDistributiveConditionalWithInfer.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.fun = void 0;
// This function's type is changed on declaration
const fun = (subFun) => { };
exports.fun = fun;


//// [declarationEmitDistributiveConditionalWithInfer.d.ts]
// This function's type is changed on declaration
export declare const fun: (subFun: <Collection, Field extends keyof Collection>() => (Collection[Field] extends infer T ? T extends Collection[Field] ? T extends readonly (infer InnerArr)[] ? infer InnerArr : T : never : never)[]) => void;


//// [DtsFileErrors]


declarationEmitDistributiveConditionalWithInfer.d.ts(2,193): error TS1338: 'infer' declarations are only permitted in the 'extends' clause of a conditional type.


==== declarationEmitDistributiveConditionalWithInfer.d.ts (1 errors) ====
    // This function's type is changed on declaration
    export declare const fun: (subFun: <Collection, Field extends keyof Collection>() => (Collection[Field] extends infer T ? T extends Collection[Field] ? T extends readonly (infer InnerArr)[] ? infer InnerArr : T : never : never)[]) => void;
                                                                                                                                                                                                    ~~~~~~~~~~~~~~
!!! error TS1338: 'infer' declarations are only permitted in the 'extends' clause of a conditional type.
    