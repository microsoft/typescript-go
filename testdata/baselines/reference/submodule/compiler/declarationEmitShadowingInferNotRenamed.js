//// [tests/cases/compiler/declarationEmitShadowingInferNotRenamed.ts] ////

//// [declarationEmitShadowingInferNotRenamed.ts]
// Any instance type
type Client = string

// Modified instance
type UpdatedClient<C> = C & {foo: number}

export const createClient = <
  D extends
    | (new (...args: any[]) => Client) // accept class
    | Record<string, new (...args: any[]) => Client> // or map of classes
>(
  clientDef: D
): D extends new (...args: any[]) => infer C
  ? UpdatedClient<C> // return instance
  : {
      [K in keyof D]: D[K] extends new (...args: any[]) => infer C // or map of instances respectively
        ? UpdatedClient<C>
        : never
    } => {
  return null as any
}

//// [declarationEmitShadowingInferNotRenamed.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createClient = void 0;
const createClient = (clientDef) => {
    return null;
};
exports.createClient = createClient;


//// [declarationEmitShadowingInferNotRenamed.d.ts]
// Modified instance
type UpdatedClient<C> = C & {
    foo: number;
};
export declare const createClient: <D extends Record<string, new (...args: any[]) => string> | (new (...args: any[]) => string)>(clientDef: D) => D extends new (...args: any[]) => infer C ? UpdatedClient<infer C> : { [K in keyof D]: D[K] extends new (...args: any[]) => infer C // or map of instances respectively
 ? UpdatedClient<infer C // or map of instances respectively
> : never; };
export {};


//// [DtsFileErrors]


declarationEmitShadowingInferNotRenamed.d.ts(5,205): error TS1338: 'infer' declarations are only permitted in the 'extends' clause of a conditional type.
declarationEmitShadowingInferNotRenamed.d.ts(6,18): error TS1338: 'infer' declarations are only permitted in the 'extends' clause of a conditional type.


==== declarationEmitShadowingInferNotRenamed.d.ts (2 errors) ====
    // Modified instance
    type UpdatedClient<C> = C & {
        foo: number;
    };
    export declare const createClient: <D extends Record<string, new (...args: any[]) => string> | (new (...args: any[]) => string)>(clientDef: D) => D extends new (...args: any[]) => infer C ? UpdatedClient<infer C> : { [K in keyof D]: D[K] extends new (...args: any[]) => infer C // or map of instances respectively
                                                                                                                                                                                                                ~~~~~~~
!!! error TS1338: 'infer' declarations are only permitted in the 'extends' clause of a conditional type.
     ? UpdatedClient<infer C // or map of instances respectively
                     ~~~~~~~
!!! error TS1338: 'infer' declarations are only permitted in the 'extends' clause of a conditional type.
    > : never; };
    export {};
    